package webinterface

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var backlogLimit int

var wsLogHook *WSLogHook

func NewWSLogHook() *WSLogHook {
	wsLogHook = new(WSLogHook)
	backlogLimit = viper.GetInt("interface.backloglimit")
	fmt.Println(backlogLimit)
	return wsLogHook
}

func MarshalEntry(e *logrus.Entry) ([]byte, error) {
	t := e.Time.Format("2006-01-02 15:04:05")

	fmt.Println(e.Data)

	return json.Marshal(WSLogEntry{Level: e.Level.String(), Message: e.Message, AtTime: t})
}

func (h *WSLogHook) Fire(e *logrus.Entry) error {
	h.AddLogItem(e)
	Connections.SendEntry(h.entries[len(h.entries)-1])
	return nil
}

func (h *WSLogHook) AddLogItem(e *logrus.Entry) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if json, err := MarshalEntry(e); err == nil {
		h.entries = append(h.entries, json)
	}

	if l := len(h.entries); l > backlogLimit {
		h.entries = h.entries[l-backlogLimit:]
	}
	return nil
}

func (h *WSLogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *WSLogHook) SendLog(conn *WsConnection) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, message := range h.entries {

		fmt.Println(string(message))

		conn.SendJson(message)
	}
}
