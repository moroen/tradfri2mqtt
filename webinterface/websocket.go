package webinterface

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
	"github.com/rs/xid"
)

type WSLogHook struct {
	mu      sync.Mutex
	entries [][]byte
}

type WSLogEntry struct {
	Level   string `json:"level"`
	Message string `json:"message"`
	AtTime  string `json:"time"`
}

var Connections WsConnections

type WsConnections struct {
	mu          sync.Mutex
	connections []*WsConnection
}

type WsConnection struct {
	Id            string
	Connection    *websocket.Conn
	ShouldSendLog bool
}

func (c *WsConnections) Add(conn *websocket.Conn) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	nc := WsConnection{}
	nc.Init(conn)
	c.connections = append(c.connections, &nc)
	return nil
}

func (c *WsConnections) Remove(Id string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	for i := range c.connections {
		if c.connections[i].Id == Id {
			c.connections[i] = c.connections[len(c.connections)-1]
			c.connections = c.connections[:len(c.connections)-1]
			logrus.WithFields(logrus.Fields{
				"Id": Id,
			}).Debug("Connection removed")
			return nil
		}
	}
	return nil
}

func (c *WsConnections) SendEntry(message []byte) {

	for _, conn := range c.connections {
		if conn.ShouldSendLog {
			fmt.Println("SendEntry", conn.Id)
			conn.SendJson(message)
		}
	}

}

func (c *WsConnection) Close(code int, text string) error {
	logrus.WithFields(logrus.Fields{
		"Id": c.Id,
	}).Info("Closing connection")
	Connections.Remove(c.Id)
	return nil
}

func (c *WsConnection) Init(conn *websocket.Conn) error {
	c.Id = xid.New().String()
	c.Connection = conn
	c.Connection.SetCloseHandler(c.Close)
	c.Read()
	logrus.WithFields(logrus.Fields{
		"Id": c.Id,
	}).Info("New connection")
	return nil
}

func (c *WsConnection) SendJson(json []byte) error {
	c.Connection.SetWriteDeadline(time.Now().Add(time.Second * 5))
	if err := c.Connection.WriteMessage(websocket.TextMessage, json); err != nil {
		fmt.Println(err.Error())
		return err
	} else {
		return err
	}
}

func (c *WsConnection) Read() {
	go func() {
		for {
			_, msg, err := c.Connection.ReadMessage()
			if err != nil {
				fmt.Println(err.Error())
				return
			} else {
				logrus.WithFields(logrus.Fields{
					"Id":      c.Id,
					"Message": string(msg),
				}).Debug("Connection read")
				var cmd WSocketCommand
				if err := json.Unmarshal(msg, &cmd); err == nil {
					switch cmd.Class {
					case "log":
						switch cmd.Command {
						case "start":
							wsLogHook.SendLog(c)

							c.ShouldSendLog = true
						case "stop":
							c.ShouldSendLog = false
						}

					}
				}

			}
		}
	}()
}

/*
func SendJson(json []byte) error {
	fmt.Println("SendJSon - Sending message")
	if wsConnection != nil {
		if err := wsConnection.WriteMessage(websocket.TextMessage, json); err != nil {
			fmt.Println("SendJson - Error", err.Error())
			return err
		} else {
			fmt.Println("SendJson - Message sent")
			return err
		}
	} else {
		return errors.New("sendJson - Not connected")
	}
}
*/

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func wshandler(w http.ResponseWriter, r *http.Request) {
	if wsConnection, err := wsupgrader.Upgrade(w, r, nil); err != nil {
		fmt.Printf("Failed to set websocket upgrade: %+v\n", err)
		return
	} else {
		Connections.Add(wsConnection)
	}
}

/*
func wshandler(w http.ResponseWriter, r *http.Request) {
	var err error

	if wsConnection, err = wsupgrader.Upgrade(w, r, nil); err != nil {
		fmt.Printf("Failed to set websocket upgrade: %+v\n", err)
		return
	} else {
		_wsViper.Set("sendLog", false)
		// wsConnection.SetWriteDeadline(time.Now().Add(time.Second * 10))
		wsConnection.SetCloseHandler(func(code int, text string) error {
			_wsViper.Set("sendLog", false)
			fmt.Println("close")
			fmt.Println(code, text)
			wsConnection = nil
			return nil
		})
		logrus.Info("Websocket: Connected")
		go func() {
			for {
				_, msg, err := wsConnection.ReadMessage()
				if err != nil {
					fmt.Println("test")
					fmt.Println(err.Error())

					return
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
*/

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
