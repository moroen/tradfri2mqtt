package webinterface

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
)

var wsLogHook *WSLogHook

type WSLogEntry struct {
	Level   string `json:"level"`
	Message string `json:"message"`
	AtTime  string `json:"time"`
}

func NewWSLogHook() *WSLogHook {
	wsLogHook = new(WSLogHook)
	// backlogLimit = viper.GetInt("interface.backloglimit")
	return wsLogHook
}

func MarshalEntry(e *logrus.Entry) ([]byte, error) {
	t := e.Time.Format("2006-01-02 15:04:05")
	entry := WsMessage{Class: "log"}
	entry.Data = WSLogEntry{Level: e.Level.String(), Message: e.Message, AtTime: t}
	return json.Marshal(entry)
}

func (h *WSLogHook) Fire(e *logrus.Entry) error {

	if msg, err := MarshalEntry(e); err == nil {
		if hub == nil {
			fmt.Println(msg)
			return ErrorWebSocketHubNotDefined
		}
		select {
		case hub.broadcast <- msg:
		default:

		}
	}
	return nil
}

func (h *WSLogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
