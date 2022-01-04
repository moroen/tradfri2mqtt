package mqttclient

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/moroen/tradfri2mqtt/tradfri"
	log "github.com/sirupsen/logrus"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTpayload struct {
	Value int64
}

type MQTTStatePayload struct {
	State string
}

type MQTTboolPayload struct {
	Value string
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
		fmt.Println(ItemQueue)
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

	if token := client.Subscribe("tradfri/+/WS/set", 0, SetHex); token.Wait() && token.Error() != nil {
		log.Print(token.Error())
	}

	if token := client.Subscribe("tradfri/+/CWS/set", 0, SetHex); token.Wait() && token.Error() != nil {
		log.Print(token.Error())
	}

	if token := client.Subscribe("tradfri/+/dimmer/set", 0, Dimmer); token.Wait() && token.Error() != nil {
		log.Print(token.Error())
	}

	if token := client.Subscribe("tradfri/+/blind/set", 0, Blind); token.Wait() && token.Error() != nil {
		log.Print(token.Error())
	}

	if token := client.Subscribe("tradfri/+/switch/set", 0, State); token.Wait() && token.Error() != nil {
		log.Print(token.Error())
	}

	if token := client.Subscribe("tradfri/cmd/#", 0, Command); token.Wait() && token.Error() != nil {
		log.Print(token.Error())
	}

}

type MQTTStateMessage struct {
	Color struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	} `json:"color"`
	State      string `json:"state"`
	Brightness int    `json:"brightness"`
	ColorTemp  int    `json:"color_temp"`
	Position   int    `json:"position"`
}

func ParseMessage(msg mqtt.Message) (deviceid int64, state int, level int, x int64, y int64, color_temp int, err error) {

	message := MQTTStateMessage{Brightness: -1, Position: -1}
	message.ColorTemp = -1
	message.Color.X = -1
	message.Color.Y = -1

	if err = json.Unmarshal(msg.Payload(), &message); err == nil {
		s := strings.Split(msg.Topic(), "/")

		deviceid, err = strconv.ParseInt(s[1], 10, 64)
		if err != nil {
			log.Fatalln(err.Error())
			return 0, 0, 0, -1, -1, -1, err
		}

		switch message.State {
		case "ON":
			state = 1
		case "OFF":
			state = 0
		default:
			state = -1
		}

		if message.Color.X != -1 {
			x = int64(message.Color.X * 65535)
		} else {
			x = -1
		}

		if message.Color.Y != -1 {
			y = int64(message.Color.Y * 65535)
		} else {
			y = -1
		}

		if message.Position != -1 {
			return deviceid, state, message.Position, x, y, message.ColorTemp, nil
		} else {
			return deviceid, state, message.Brightness, x, y, message.ColorTemp, nil
		}
	} else {
		log.WithFields(log.Fields{
			"Error": err.Error(),
		}).Error("PrseMessage - Unmarshal")
		return 0, 0, 0, -1, -1, -1, err
	}
}

func State(client mqtt.Client, msg mqtt.Message) {
	state := -1

	log.WithFields(log.Fields{
		"topic":   msg.Topic(),
		"payload": msg.Payload(),
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
			if payload.Value == "true" {
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

	tradfri.State(int(deviceid), state)
}

func SetHex(client mqtt.Client, msg mqtt.Message) {
	/* 	log.WithFields(log.Fields{
	   		"topic":   msg.Topic(),
	   		"payload": string(msg.Payload()),
	   	}).Debug("Received color message")

	   	deviceid, value, _, _, _, err := ParseMessage(msg)
	   	if err != nil {
	   		log.Error(err.Error())
	   		return
	   	}

	   	if err = coap.SetHexForLevel(deviceid, value); err != nil {
	   		log.WithFields(log.Fields{
	   			"error": err.Error(),
	   		}).Error("Color message error")
	   	} */
}

func Blind(client mqtt.Client, msg mqtt.Message) {
	log.WithFields(log.Fields{
		"topic":   msg.Topic(),
		"payload": string(msg.Payload()),
	}).Debug("Received blind message")

	/*
		deviceid, _, value, _, _, _, err := ParseMessage(msg)
		if err != nil {
			log.Error(err.Error())
			return
		}

		log.WithFields(log.Fields{
			"DeviceID": deviceid,
			"Position": value,
		}).Debug("Handlers - Blind")

		if _, err := coap.SetBlind(deviceid, value); err != nil {
			log.WithFields(log.Fields{
				"Error": err.Error(),
			}).Error("Handlers - Blind")
		}
	*/
}

func Dimmer(client mqtt.Client, msg mqtt.Message) {
	// fmt.Printf("Received dimmer message: %s from topic: %s\n", msg.Payload(), msg.Topic())

	log.WithFields(log.Fields{
		"topic":   msg.Topic(),
		"payload": string(msg.Payload()),
	}).Debug("Received dimmer message")

	deviceid, state, level, x, y, col_temp, err := ParseMessage(msg)
	if err != nil {
		log.Error(err.Error())
		return
	}

	log.WithFields(log.Fields{
		"State":    state,
		"Level":    level,
		"X":        x,
		"Y":        y,
		"Col_Temp": col_temp,
	}).Debug("MQTT-handlers - Dimmer")

	if state != -1 {
		tradfri.State(int(deviceid), state)
	}

	if level != -1 {
		tradfri.Level(int(deviceid), level)
	}

	if x != -1 {
		tradfri.SetXY(int(deviceid), int(x), int(y))
	}

	if col_temp != -1 {
		log.WithFields(log.Fields{
			"col_temp": col_temp,
		}).Debug("MQTT-handlers - Dimmer - SetColorTemp")
		if col_temp >= 350 {
			tradfri.SetHex(int(deviceid), "efd275")
		} else if col_temp < 280 {
			tradfri.SetHex(int(deviceid), "f5faf6")
		} else {
			tradfri.SetHex(int(deviceid), "f1e0b5")
		}
	}

	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Dimmer message error")
	}
}
