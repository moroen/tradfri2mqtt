package mqttclient

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	coap "github.com/moroen/go-tradfricoap"
)

type MQTTmessage struct {
	Value int `json:"value"`
}

type MQTTboolMessage struct {
	Value bool `json:"value"`
}

type MQTTdimmmerMessage struct {
	State      string `json:"state"`
	Brightness int    `json:"brightness"`
}

func Show(msg []byte) error {
	// fmt.Printf("%s\n", msg)

	var topic string
	var valueJson []byte

	SendConfigObject(msg)

	if info, err := coap.ParseLightInfo(msg); err == nil {
		topic = fmt.Sprintf("tradfri/%d/38/0/dimmer", info.Id)

		var state string
		if info.State {
			state = "ON"
		} else {
			state = "OFF"
		}

		if message, err := json.Marshal(MQTTdimmmerMessage{State: state, Brightness: int(info.Dimmer)}); err == nil {
			log.WithFields(log.Fields{
				"Function": "Show - Send dimmer message",
				"Topic":    topic,
				"Message":  string(message),
			}).Debug()
			SendTopic(topic, message)
		} else {
			log.WithFields(log.Fields{
				"Function": "Show - Send dimmer message",
				"Error":    err.Error(),
			}).Error()
		}

		/*
			if info.ColorSpace != "" {
				topic = fmt.Sprintf("tradfri/%d/38/1/%s", info.Id, info.ColorSpace)
				if valueJson, err = json.Marshal(MQTTmessage{Name: info.Name, Value: info.HexLevel()}); err == nil {
					SendTopic(topic, valueJson)
				}
			}
		*/

		if err != nil {
			log.Fatal(err.Error())
		}
		return err

	} else if info, err := coap.ParsePlugInfo(msg); err == nil {
		topic = fmt.Sprintf("tradfri/%d/37/0/switch", info.Id)
		if valueJson, err = json.Marshal(MQTTboolMessage{Value: info.State}); err == nil {
			return SendTopic(topic, valueJson)
		} else {
			log.Error(err.Error())
			return err
		}

	} else if info, err := coap.ParseBlindInfo(msg); err == nil {
		topic = fmt.Sprintf("tradfri/%d/38/0/blind", info.Id)
		if valueJson, err = json.Marshal(MQTTmessage{Value: int(info.Level)}); err == nil {
			return SendTopic(topic, valueJson)
		} else {
			log.Error(err.Error())
			return err
		}
	} else {
		/*
			log.WithFields(log.Fields{
				"Coap message": string(msg),
				"Error":        err.Error(),
			}).Info()
		*/
		return err
	}
}
