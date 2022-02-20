package webinterface

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
)

type WSLogHook struct {
	mu      sync.Mutex
	entries [][]byte
}

type WSLogEntry struct {
	Level   string    `json:"level"`
	Message string    `json:"message"`
	AtTime  time.Time `json:"time"`
}

var wsLogHook *WSLogHook

var wsConnection *websocket.Conn
var backlogLimit int

func GetLogHook() *WSLogHook {
	wsLogHook = new(WSLogHook)
	backlogLimit = viper.GetInt("interface.backloglimit")
	fmt.Println(backlogLimit)
	return wsLogHook
}

func MarshalEntry(e *logrus.Entry) ([]byte, error) {
	return json.Marshal(WSLogEntry{Level: e.Level.String(), Message: e.Message, AtTime: e.Time})
}

func (h *WSLogHook) Fire(e *logrus.Entry) error {
	h.AddLogItem(e)
	if _wsViper.GetBool("sendLog") {
		h.SendLatestEntry()
	}
	return nil
}

func SendJson(json []byte) error {
	if wsConnection != nil {
		if err := wsConnection.WriteMessage(websocket.TextMessage, json); err != nil {
			logrus.WithFields(logrus.Fields{
				"Error": err.Error(),
			}).Error("WebSocket - SendJson")
			return err
		} else {
			return err
		}
	} else {
		return errors.New("sendJson - Not connected")
	}
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

func (h *WSLogHook) SendLog(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, message := range h.entries {
		SendJson(message)
	}
}

func (h *WSLogHook) SendLatestEntry() {
	h.mu.Lock()
	defer h.mu.Unlock()
	SendJson(h.entries[len(h.entries)-1])
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func wshandler(w http.ResponseWriter, r *http.Request) {
	var err error

	if wsConnection, err = wsupgrader.Upgrade(w, r, nil); err != nil {
		fmt.Printf("Failed to set websocket upgrade: %+v\n", err)
		return
	} else {
		logrus.Info("Websocket: Connected")
		go func() {
			for {
				_, msg, err := wsConnection.ReadMessage()
				if err != nil {
					break
				}
				var cmd WSocketCommand
				if err := json.Unmarshal(msg, &cmd); err == nil {
					switch cmd.Class {
					case "log":
						switch cmd.Command {
						case "start":
							wsLogHook.SendLog(wsConnection)
							_wsViper.Set("sendLog", true)
						case "stop":
							_wsViper.Set("sendLog", false)
						}
					}
				} else {
					logrus.Error("Websocket: Unable to unmarshal command")
				}
			}
		}()
	}
}

type WSocketCommand struct {
	Class   string `json:"class"`
	Command string `json:"command"`
	Value   string `json:"value"`
}

func WebSocketRoutes(r *gin.Engine) error {
	r.GET("/api/ws", func(c *gin.Context) {
		wshandler(c.Writer, c.Request)
	})

	r.POST("/api/ws", func(c *gin.Context) {

		fmt.Println("post")
		var cmd WSocketCommand

		if err := c.ShouldBind(&cmd); err == nil {
			fmt.Println(cmd)
			c.JSON(200, PostResponse{Status: "Ok"})
		}

	})

	return nil
}
