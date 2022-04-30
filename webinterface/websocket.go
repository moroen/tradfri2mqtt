package webinterface

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/moroen/tradfri2mqtt/tradfri"
	"github.com/sirupsen/logrus"
)

type WSLogHook struct {
	mu      sync.Mutex
	entries [][]byte
}

type WsMessage struct {
	Class string      `json:"class"`
	Data  interface{} `json:"data"`
}

type MQTTMessage struct {
	Topic   string
	Payload []byte
}

func HandleCommand(msg []byte) error {
	var cmd WSocketCommand
	if err := json.Unmarshal(msg, &cmd); err == nil {

		// fmt.Printf("%+v\n", cmd)
		switch cmd.Class {
		case "log":
			switch cmd.Command {
			case "start":
			case "stop":
			}
		case "devices":
			switch cmd.Command {
			case "update":
				tradfri.Discover(true)
			}
		case "device":
			msg := WSDeviceMessage{}
			json.Unmarshal(*cmd.Value, &msg)
			topic := fmt.Sprintf("internal/%d", msg.Deviceid)
			// fmt.Printf("%+v\n", msg)
			tradfri.Dimmer(topic, *msg.Payload)
		}
	} else {
		fmt.Println(err.Error())
	}
	return nil
}

func SendJson(message []byte) error {
	select {
	case hub.broadcast <- message:
	default:
	}
	return nil
}

func SendDeviceJSON_not_enabled(message interface{}) error {
	logrus.Debug("SendDeviceJSON - Interface not enabled")
	return nil
}

func SendDeviceJSON(message interface{}) error {

	if ws, err := json.Marshal(WsMessage{
		Class: "devices",
		Data:  message,
	}); err == nil {
		SendJson(ws)

	}

	return nil
}

type WSocketCommand struct {
	Class   string           `json:"class"`
	Command string           `json:"command"`
	Value   *json.RawMessage `json:"value"`
}

type WSDeviceMessage struct {
	Deviceid int64            `json:"deviceid"`
	Payload  *json.RawMessage `json:"payload"`
}

type WSPayloadState struct {
	Value bool `json:"value"`
}
