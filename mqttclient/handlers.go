package mqttclient

import (
	"encoding/json"
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

type ColorMessage struct {
	Color struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	} `json:"color"`
	State      string `json:"state"`
	Brightness int    `json:"brightness"`
}

func ParseMessage(msg mqtt.Message) (deviceid int64, state int, level int, x int64, y int64, err error) {

	message := ColorMessage{Brightness: -1}
	message.Color.X = -1
	message.Color.Y = -1

	if err = json.Unmarshal(msg.Payload(), &message); err == nil {
		pretty_print(message)
		s := strings.Split(msg.Topic(), "/")

		deviceid, err = strconv.ParseInt(s[1], 10, 64)
		if err != nil {
			log.Fatalln(err.Error())
			return 0, 0, 0, -1, -1, err
		}

		switch message.State {
		case "ON":
			state = 1
		case "OFF":
			state = 0
		default:
			state = -1
		}

		/*
			if message.Color.R != -1 {
				hex = fmt.Sprintf("#%02x%02x%02x", message.Color.R, message.Color.G, message.Color.B)
			} else {
				hex = ""
			}
		*/
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

		return deviceid, state, message.Brightness, x, y, nil
	} else {
		log.WithFields(log.Fields{
			"Function": "ParseMessage - Unmarshal",
			"Error":    err.Error(),
		})
		return 0, 0, 0, -1, -1, err
	}
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

	deviceid, _, value, _, _, err := ParseMessage(msg)
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

	deviceid, state, level, x, y, err := ParseMessage(msg)
	if err != nil {
		log.Error(err.Error())
		return
	}

	log.WithFields(log.Fields{
		"State": state,
		"Level": level,
		"X":     x,
		"Y":     y,
	}).Debug()

	if state != -1 {
		if err := coap.SetState(deviceid, state); err != nil {
			log.WithFields(log.Fields{
				"Function": "MQTT-handlers - Dimmer - SetState",
				"Error":    err.Error(),
			}).Error()
		}
	}

	if level != -1 {
		if err := coap.SetLevel(deviceid, level); err != nil {
			log.WithFields(log.Fields{
				"Function": "MQTT-handlers - Dimmer - SetLevel",
				"Error":    err.Error(),
			}).Error()
		}
	}

	if x != -1 {
		if err := coap.SetXY(deviceid, x, y); err != nil {
			log.WithFields(log.Fields{
				"Function": "MQTT-handlers - Dimmer - SetXY",
				"Error":    err.Error(),
			}).Error()
		}
	}

	/*
		if rgb != "" {
			if err := coap.SetRGB(deviceid, rgb); err != nil {
				log.WithFields(log.Fields{
					"Function": "MQTT-handlers - Dimmer - SetColor",
					"Error":    err.Error(),
					"RGB":      rgb,
				}).Error()
			}
		}
	*/

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
