package mqttclient

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/moroen/tradfri2mqtt/settings"
	log "github.com/sirupsen/logrus"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var _client mqtt.Client
var _status_channel chan (error)

var _wg *sync.WaitGroup

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	opts := client.OptionsReader()
	log.Info(fmt.Sprintf("Connected to broker at %s", opts.Servers()))
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Info(fmt.Sprintf("Connect lost: %v", err))
	connectToBroker(_client)
}

func GetMQTTClientConnection() mqtt.Client {
	return _client
}

func SendTopic(topic string, payload []byte) error {
	if client := GetMQTTClientConnection(); client != nil {
		client.Publish(topic, 0, false, payload)
	} else {
		log.WithFields(log.Fields{
			"topic":   string(topic),
			"payload": string(payload),
		}).Error("MQTT client connection not set")
	}
	return nil
}

func connectToBroker(client mqtt.Client) error {
	retries := 1
	waiting := 10

	for retries < 100 {
		if token := client.Connect(); token.Wait() && token.Error() == nil {
			Subscribe(client, _status_channel)
			return nil
		} else {
			log.Debug(fmt.Sprintf("Unable to connect to broker. Waiting %d seconds before trying again", waiting))
			time.Sleep(time.Duration(waiting) * time.Second)
			retries++
			if waiting = waiting * retries; waiting > 20 {
				waiting = 20
			}

		}
	}
	return nil
}

func Start(wg *sync.WaitGroup, status_channel chan (error)) {
	_status_channel = status_channel
	_wg = wg

	_wg.Add(1)

	log.Info("MQTT: Starting")

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

	if err := connectToBroker(_client); err == nil {
		go HandleQueue()
	} else {
		status_channel <- errors.New("Broker not connected, giving up")
	}

}

func Stop() {
	if _wg != nil {
		defer _wg.Done()
	}

	if _client != nil {
		log.Info("MQTT: Stopping")
		_client.Disconnect(250)
		log.Info("MQTT: Stopped")

	}
}

func Restart() {
	Stop()
	go Start(_wg, _status_channel)
}

func Do_Test(wg *sync.WaitGroup) {
	defer wg.Done()
	wg.Add(1)
}
