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
	ColorTemp  int    `json:"color_temp"`
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

		var color_temp int

		if info.ColorSpace == "WS" {
			colorMap := coap.CWmap()
			switch info.Hex {
			case colorMap[10]["Hex"]:
				color_temp = 50
			case colorMap[20]["Hex"]:
				color_temp = 320
			case colorMap[30]["Hex"]:
				color_temp = 480
			}
		}

		message := MQTTStateMessage{State: state, Brightness: int(info.Dimmer), ColorTemp: color_temp}
		message.Color.X, message.Color.Y = info.ColorValues.ToFloat()

		if messageJson, err := json.Marshal(message); err == nil {
			log.WithFields(log.Fields{
				"Topic":   topic,
				"Message": string(messageJson),
			}).Debug("Show - Send dimmer message")
			SendTopic(topic, messageJson, false)
		} else {
			log.WithFields(log.Fields{
				"Error": err.Error(),
			}).Error("Show - Send dimmer message")
		}

		if err != nil {
			log.Fatal(err.Error())
		}
		return err

	} else if info, err := coap.ParsePlugInfo(msg); err == nil {
		topic = fmt.Sprintf("tradfri/%d/37/0/switch", info.Id)
		if valueJson, err = json.Marshal(MQTTboolMessage{Value: info.State}); err == nil {
			return SendTopic(topic, valueJson, false)
		} else {
			log.Error(err.Error())
			return err
		}

	} else if info, err := coap.ParseBlindInfo(msg); err == nil {
		topic = fmt.Sprintf("tradfri/%d/blind", info.Id)
		if valueJson, err = json.Marshal(MQTTmessage{Value: int(info.Level)}); err == nil {
			return SendTopic(topic, valueJson, false)
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
