package messages

import (
	"encoding/json"
	"fmt"
	"math"

	log "github.com/sirupsen/logrus"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	coap "github.com/moroen/go-tradfricoap"
)

var _client mqtt.Client

type MQTTmessage struct {
	Value int `json:"value"`
}

type MQTTboolMessage struct {
	Value bool   `json:"value"`
	Name  string `json:"nodeName"`
}

func SetClientConnection(client mqtt.Client) {
	_client = client
}

func SendTopic(topic string, payload []byte) error {
	if _client != nil {
		if topic != "" {
			_client.Publish(topic, 0, false, payload)
		}
	} else {
		log.Fatal("MQTT client connection not set")
	}
	return nil
}

func Show(msg []byte) error {
	// fmt.Printf("%s\n", msg)

	var topic string
	var valueJson []byte

	if info, err := coap.ParseLightInfo(msg); err == nil {
		topic = fmt.Sprintf("tradfri/%d/38/0/dimmer", info.Id)
		if info.State {
			valueJson, err = json.Marshal(MQTTmessage{Value: int(math.Round(float64(info.Dimmer) / 2.54))})
		} else {
			valueJson, err = json.Marshal(MQTTmessage{Value: 0})
		}

		SendTopic(topic, valueJson)

		if info.ColorSpace != "" {
			topic = fmt.Sprintf("tradfri/%d/38/1/%s", info.Id, info.ColorSpace)
			if valueJson, err = json.Marshal(MQTTmessage{Value: info.HexLevel()}); err == nil {
				SendTopic(topic, valueJson)
			}
		}

		if err != nil {
			log.Println(err.Error())
		}
		return err

	} else if info, err := coap.ParsePlugInfo(msg); err == nil {
		topic = fmt.Sprintf("tradfri/%d/37/0/switch", info.Id)
		valueJson, err = json.Marshal(MQTTboolMessage{Value: info.State, Name: info.Name})
		if err != nil {
			log.Println(err.Error())
			return err
		}
	} else if info, err := coap.ParseBlindInfo(msg); err == nil {
		topic = fmt.Sprintf("tradfri/%d/38/0/blind", info.Id)
		valueJson, err = json.Marshal(MQTTmessage{Value: int(info.Level)})
		if err != nil {
			log.Println(err.Error())
			return err
		}
	} else {
		log.Println(err.Error())
	}

	return SendTopic(topic, valueJson)
}
