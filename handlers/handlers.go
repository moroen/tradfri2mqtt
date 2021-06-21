package handlers

import (
	"encoding/json"
	"math"
	"strconv"
	"strings"

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

func Subscribe(client mqtt.Client) {
	if token := client.Subscribe("tradfri/+/38/+/cw/set", 0, SetHex); token.Wait() && token.Error() != nil {
		log.Print(token.Error())
	}

	if token := client.Subscribe("tradfri/+/38/+/cws/set", 0, SetHex); token.Wait() && token.Error() != nil {
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

	if _, err = coap.SetHexForLevel(deviceid, value); err != nil {
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
		_, err = coap.SetState(deviceid, 1)
	case 0:
		_, err = coap.SetState(deviceid, 0)
	default:
		_, err = coap.SetLevel(deviceid, int(math.Round(float64(value)*2.54)))
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

	_, err = coap.SetState(deviceid, state)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Dimmer message error")
	}
}
