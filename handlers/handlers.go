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
	if token := client.Subscribe("tradfri/+/38/+/dimmer/set", 0, Dimmer); token.Wait() && token.Error() != nil {
		log.Print(token.Error())
	}

	if token := client.Subscribe("tradfri/+/37/+/switch/set", 0, State); token.Wait() && token.Error() != nil {
		log.Print(token.Error())
	}

	if token := client.Subscribe("tradfri/cmd/#", 0, Command); token.Wait() && token.Error() != nil {
		log.Print(token.Error())
	}

}

func Dimmer(client mqtt.Client, msg mqtt.Message) {
	// fmt.Printf("Received dimmer message: %s from topic: %s\n", msg.Payload(), msg.Topic())

	log.WithFields(log.Fields{
		"topic":   msg.Topic(),
		"payload": string(msg.Payload()),
	}).Debug("Received dimmer message")

	var payload MQTTpayload
	err := json.Unmarshal(msg.Payload(), &payload)
	if err != nil {
		log.Fatalln(err.Error())
	}

	s := strings.Split(msg.Topic(), "/")

	deviceid, err := strconv.ParseInt(s[1], 10, 64)
	if err != nil {
		log.Fatalln(err.Error())
	}

	if payload.Value == 255 {
		err = coap.SetState(deviceid, 1)
	} else if payload.Value == 0 {
		err = coap.SetState(deviceid, 0)
	} else {
		err = coap.SetLevel(deviceid, int(math.Round(float64(payload.Value)*2.54)))
	}

	if err != nil {
		log.WithFields(log.Fields{
			"topic":   msg.Topic(),
			"payload": string(msg.Payload()),
			"error":   err.Error(),
			// "payload": string(msg.Payload()),
		}).Error("Dimmer message failed")
	} else {
		log.WithFields(log.Fields{
			"topic":   msg.Topic(),
			"payload": string(msg.Payload()),
		}).Debug("Dimmer message done")
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
		log.WithFields(log.Fields{
			"topic": msg.Topic(),
			"error": err.Error(),
			// "payload": string(msg.Payload()),
		}).Error("State message failed")
	} else {
		log.WithFields(log.Fields{
			"topic":   msg.Topic(),
			"payload": string(msg.Payload()),
		}).Debug("State message done")
	}
}
