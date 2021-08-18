package mqttclient

import (
	"encoding/json"
	"math"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	coap "github.com/moroen/go-tradfricoap"
)

type MQTTpayload struct {
	Value int
}

type MQTTboolPayload struct {
	Value bool
}

type MQTTgwConfig struct {
	Ip  string
	Key string
}

type QueueItem struct {
	Client  mqtt.Client
	Message mqtt.Message
	Handler func(mqtt.Client, mqtt.Message)
	Retries int
}

var ItemQueue []QueueItem

func AddToQueue(client mqtt.Client, msg mqtt.Message, handler func(mqtt.Client, mqtt.Message)) {
	var item = QueueItem{Client: client, Message: msg, Handler: handler}
	ItemQueue = append(ItemQueue, item)
}

func HandleQueue() {
	for {
		if len(ItemQueue) > 0 {
			log.Debug("Handling queue")
			ItemQueue[0].Handler(ItemQueue[0].Client, ItemQueue[0].Message)
			ItemQueue = ItemQueue[1:]
		}
		time.Sleep(time.Second * 2)
	}
}

func Subscribe(client mqtt.Client, status_channel chan (error)) {
	_status_channel = status_channel

	if client == nil {
		log.Fatal("Unable to subscribe, client is nil")
		return
	}

	if token := client.Subscribe("tradfri/+/38/+/WS/set", 0, SetHex); token.Wait() && token.Error() != nil {
		log.Print(token.Error())
	}

	if token := client.Subscribe("tradfri/+/38/+/CWS/set", 0, SetHex); token.Wait() && token.Error() != nil {
		log.Print(token.Error())
	}

	if token := client.Subscribe("tradfri/+/38/+/dimmer/set", 0, Dimmer); token.Wait() && token.Error() != nil {
		log.Print(token.Error())
	}

	if token := client.Subscribe("tradfri/+/38/+/blind/set", 0, Blind); token.Wait() && token.Error() != nil {
		log.Print(token.Error())
	}

	if token := client.Subscribe("tradfri/+/37/+/switch/set", 0, State); token.Wait() && token.Error() != nil {
		log.Print(token.Error())
	}

	if token := client.Subscribe("tradfri/cmd/#", 0, Command); token.Wait() && token.Error() != nil {
		log.Print(token.Error())
	}

}

func ParseMessage(msg mqtt.Message) (deviceid int64, value int, err error) {
	var payload MQTTpayload
	err = json.Unmarshal(msg.Payload(), &payload)
	if err != nil {
		log.Fatalln(err.Error())
		return 0, 0, err
	}

	s := strings.Split(msg.Topic(), "/")

	deviceid, err = strconv.ParseInt(s[1], 10, 64)
	if err != nil {
		log.Fatalln(err.Error())
		return 0, 0, err
	}
	return deviceid, payload.Value, nil
}

func SetHex(client mqtt.Client, msg mqtt.Message) {
	log.WithFields(log.Fields{
		"topic":   msg.Topic(),
		"payload": string(msg.Payload()),
	}).Debug("Received color message")

	deviceid, value, err := ParseMessage(msg)
	if err != nil {
		log.Error(err.Error())
		return
	}

	if err = coap.SetHexForLevel(deviceid, value); err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Color message error")
	}
}

func Blind(client mqtt.Client, msg mqtt.Message) {
	log.WithFields(log.Fields{
		"topic":   msg.Topic(),
		"payload": string(msg.Payload()),
	}).Debug("Received blind message")

	deviceid, value, err := ParseMessage(msg)
	if err != nil {
		log.Error(err.Error())
		return
	}

	switch value {
	case 255:
		_, err = coap.SetBlind(deviceid, 100)
	case 0:
		_, err = coap.SetBlind(deviceid, 0)
	default:
		_, err = coap.SetBlind(deviceid, value)
	}

	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Blind message error")
	}
}

func Dimmer(client mqtt.Client, msg mqtt.Message) {
	// fmt.Printf("Received dimmer message: %s from topic: %s\n", msg.Payload(), msg.Topic())

	log.WithFields(log.Fields{
		"topic":   msg.Topic(),
		"payload": string(msg.Payload()),
	}).Debug("Received dimmer message")

	deviceid, value, err := ParseMessage(msg)
	if err != nil {
		log.Error(err.Error())
		return
	}

	switch value {
	case 255:
		err = coap.SetState(deviceid, 1)
	case 0:
		err = coap.SetState(deviceid, 0)
	default:
		err = coap.SetLevel(deviceid, int(math.Round(float64(value)*2.54)))
	}

	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Dimmer message error")
	}

}

func State(client mqtt.Client, msg mqtt.Message) {
	// fmt.Printf("Received state message: %s from topic: %s\n", msg.Payload(), msg.Topic())

	log.WithFields(log.Fields{
		"topic":   msg.Topic(),
		"payload": string(msg.Payload()),
	}).Debug("Received state message")

	var payload MQTTboolPayload
	err := json.Unmarshal(msg.Payload(), &payload)
	if err != nil {
		log.Fatalln(err.Error())
	}

	s := strings.Split(msg.Topic(), "/")

	deviceid, err := strconv.ParseInt(s[1], 10, 64)
	if err != nil {
		log.Fatalln(err.Error())
	}

	var state int
	if payload.Value {
		state = 1
	}

	err = coap.SetState(deviceid, state)
	if err != nil {
		// _status_channel <- fmt.Errorf("Plug failed")
		AddToQueue(client, msg, State)
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Dimmer message error")
	}
}
