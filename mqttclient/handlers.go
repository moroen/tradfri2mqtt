package mqttclient

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	coap "github.com/moroen/go-tradfricoap"
)

type MQTTpayload struct {
	Value int64
}

type MQTTStatePayload struct {
	State string
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

func ParseMessage(msg mqtt.Message) (deviceid int64, value int64, err error) {
	var payload MQTTpayload

	u := map[string]interface{}{}

	if err = json.Unmarshal(msg.Payload(), &u); err == nil {
		fmt.Println(u)

		if u["state"] != nil {
			if u["state"].(string) == "OFF" {
				payload = MQTTpayload{Value: 0}

			} else if u["state"].(string) == "ON" {
				payload = MQTTpayload{Value: 255}
			}
		} else if u["value"] != nil {
			payload = MQTTpayload{Value: int64(u["value"].(float64))}
		}

	} else if val, err := strconv.ParseInt(string(msg.Payload()), 10, 64); err == nil {
		if val == 99 {
			payload = MQTTpayload{Value: 255}
		} else {
			payload = MQTTpayload{Value: val}
		}

	} else {
		str := strings.ToLower(string(msg.Payload()))
		fmt.Println(str)
		if str == "off" {
			payload = MQTTpayload{Value: 0}
		} else if str == "on" {
			payload = MQTTpayload{Value: 255}
		} else {
			log.Fatalln(err.Error())
			return 0, 0, err
		}
	}

	fmt.Println(payload.Value)

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

	state := -1

	log.WithFields(log.Fields{
		"topic":         msg.Topic(),
		"payload_bytes": msg.Payload(),
		"payload":       string(msg.Payload()),
	}).Debug("Received state message")

	if bValue, err := strconv.ParseBool(string(msg.Payload())); err == nil {
		if bValue {
			state = 1
		} else {
			state = 0
		}
	} else {
		var payload MQTTboolPayload
		if err := json.Unmarshal(msg.Payload(), &payload); err == nil {
			if payload.Value {
				state = 1
			} else {
				state = 0
			}
		}
	}

	s := strings.Split(msg.Topic(), "/")

	deviceid, err := strconv.ParseInt(s[1], 10, 64)
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = coap.SetState(deviceid, state)
	if err != nil {
		// _status_channel <- fmt.Errorf("Plug failed")
		AddToQueue(client, msg, State)
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Switch message error")
	}
}
