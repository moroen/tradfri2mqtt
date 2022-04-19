package webinterface

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/moroen/tradfri2mqtt/tradfri"
)

type WSLogHook struct {
	mu      sync.Mutex
	entries [][]byte
}

type WsMessage struct {
	Class string      `json:"class"`
	Data  interface{} `json:"data"`
}

func HandleCommand(msg []byte) error {
	var cmd WSocketCommand
	if err := json.Unmarshal(msg, &cmd); err == nil {
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
			switch cmd.Command {
			case "setxy":
				fmt.Println(cmd.Value)
				tradfri.Discover(true)
			}
		}
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
	Class   string      `json:"class"`
	Command string      `json:"command"`
	Value   interface{} `json:"value"`
}
