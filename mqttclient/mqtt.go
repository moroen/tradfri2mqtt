package mqttclient

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/moroen/tradfri2mqtt/settings"
	log "github.com/sirupsen/logrus"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var ErrorBrokerConnectionRefused = errors.New("Bronker Error: Connection refused")

var _client mqtt.Client
var _status_channel chan (error)

var _wg *sync.WaitGroup

var _ctxMQTT context.Context
var _mqttStop func()

var isStopping bool

type MQTTQueueItem struct {
	Topic    string
	Payload  []byte
	Retained bool
}

var messageQueue []MQTTQueueItem

func HandleMQTTQueue() {
	var item MQTTQueueItem

	for len(messageQueue) > 0 {
		item, messageQueue = messageQueue[0], messageQueue[1:]
		SendTopic(item.Topic, item.Payload, item.Retained)
	}
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	opts := client.OptionsReader()
	log.Info(fmt.Sprintf("MQTT: Connected to broker at %s", opts.Servers()))
	Subscribe(client, _status_channel)
	go sendKeepAlive()
	go HandleMQTTQueue()
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Info(fmt.Sprintf("MQTT: Connection lost: %v", err))
	if !isStopping {
		doConnect()
	}
}

func GetMQTTClientConnection() mqtt.Client {
	return _client
}

func SendTopic(topic string, payload []byte, retained bool) error {
	if client := GetMQTTClientConnection(); client != nil {
		log.WithFields(log.Fields{
			"Topic":   topic,
			"Payload": string(payload),
		}).Debug("SendTopic")
		if token := client.Publish(topic, 0, retained, payload); token.Error() != nil {
			log.Error("Unable to publish to broker")
			messageQueue = append(messageQueue, MQTTQueueItem{Topic: topic, Payload: payload, Retained: retained})
			return token.Error()
		}

	} else {
		log.WithFields(log.Fields{
			"topic":   string(topic),
			"payload": string(payload),
		}).Error("MQTT client connection not set")
	}
	return nil
}

func connectToBroker(client mqtt.Client) error {
	if token := client.Connect(); token.Wait() && token.Error() == nil {
		return nil
	} else {
		return ErrorBrokerConnectionRefused
	}
}

func sendKeepAlive() {
	ticker := time.NewTicker(5 * time.Second)

	for {
		select {
		case <-ticker.C:
			SendTopic("tradfri/status", []byte("Alive"), false)
			break
		case <-_ctxMQTT.Done():
			SendTopic("tradfri/status", []byte("Stopping"), false)
			return
		}
	}
}

func Start(wg *sync.WaitGroup, status_channel chan (error)) {
	_status_channel = status_channel
	_wg = wg

	_wg.Add(1)

	log.Info("MQTT: Starting")

	_ctxMQTT, _mqttStop = context.WithCancel(context.Background())

	conf := settings.GetConfig(false)

	var broker = conf.Mqtt.Host
	var port = conf.Mqtt.Port
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%s", broker, port))
	opts.SetClientID("go_mqtt_client")
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	_client = mqtt.NewClient(opts)
	discovered = make(map[int64]struct{})

	doConnect()

}

func doConnect() {
	ticker := time.NewTicker(5 * time.Second)
	conf := settings.GetConfig(false)
	for {
		if err := connectToBroker(_client); err == nil {

			return
		} else {
			log.Error(fmt.Sprintf("MQTT: Unable to connect to broker at %s:%s", conf.Mqtt.Host, conf.Mqtt.Port))
		}
		select {
		case <-ticker.C:
			break
		case <-_ctxMQTT.Done():
			return
		}
	}
}

func Stop() {
	defer _wg.Done()

	isStopping = true

	_mqttStop()

	if _client != nil {
		log.Info("MQTT: Stopping...")
		_client.Disconnect(250)
		log.Info("MQTT: Stopped")
	}
}

func Restart() {
	Stop()
	go Start(_wg, _status_channel)
}
