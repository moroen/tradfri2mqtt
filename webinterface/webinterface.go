package webinterface

import (
	"bytes"
	"fmt"

	"github.com/ghodss/yaml"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/moroen/tradfri2mqtt/settings"
	"github.com/moroen/tradfri2mqtt/tradfri"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var _status_channel chan (error)
var _wsViper *viper.Viper

var r *gin.Engine

type PostResponse struct {
	Status string
	Error  string
}

type TradfriFields struct {
	Gateway string `json:"gateway"`
	Key     string `json:"key"`
}

type TradfriPSk struct {
	Identity string `json:"identity"`
	Passkey  string `json:"passkey"`
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

var _server_root string

func SetInterfaceSettingsDefaults() {
	_wsViper = viper.New()
	_wsViper.SetDefault("sendLog", false)
}

func Interface_Server(server_root string, port int, status_channel chan (error)) {

	_status_channel = status_channel
	_server_root = server_root

	SetInterfaceSettingsDefaults()

	r = gin.Default()
	r.Use(CORS())

	r.Use(static.Serve("/", static.LocalFile(_server_root, false)))

	if err := WebSocketRoutes(r); err != nil {
		log.Error("Unable to add WebSocket Routes")
	}

	if err := DeviceRoutes(r); err != nil {
		log.Error("Unable to add Device Routes")
	}

	r.GET("/api/v1/hello", func(c *gin.Context) {
		c.JSON(200, `{"message":"hello, hello, hello"}`)
		log.Info("Hello World")
	})

	r.GET("/api/settings", func(c *gin.Context) {
		// conf := settings.GetConfig(false)

		conf := viper.AllSettings()

		c.JSON(200, conf)
	})

	r.POST("/api/settings", func(c *gin.Context) {
		// conf := settings.GetConfig(false)

		/*
			if jsonData, err := ioutil.ReadAll(c.Request.Body); err == nil {
				fmt.Println(string(jsonData))
			} else {
				fmt.Println(err.Error())
			}
		*/

		var conf settings.Config

		if err := viper.Unmarshal(&conf); err != nil {
			fmt.Println(err.Error())
		}

		if err := c.ShouldBind(&conf); err == nil {
			// fmt.Printf("%+v\n", conf)

			yml, err := yaml.Marshal(conf)
			if err != nil {
				c.JSON(401, PostResponse{Status: "Error", Error: err.Error()})
				return
			}

			reader := bytes.NewReader(yml)

			if err := viper.MergeConfig(reader); err != nil {
				c.JSON(400, PostResponse{Status: "Error", Error: err.Error()})
				return
			}

			if err := viper.WriteConfig(); err != nil {
				resp := PostResponse{Status: "Error", Error: err.Error()}
				c.JSON(400, resp)
				return
			}
			c.JSON(200, PostResponse{Status: "Ok"})
			err = settings.ErrConfigIsDirty
			_status_channel <- err
		} else {
			c.JSON(400, err.Error())
		}
	})

	r.POST("/api/getPSK", func(c *gin.Context) {
		var data TradfriFields

		if err := c.ShouldBind(&data); err == nil {

			tradfri.GetNewPSK(data.Gateway, data.Key, func(ident string, key string, err error) {
				if err == nil {
					ret := TradfriPSk{Identity: ident, Passkey: key}
					c.JSON(200, ret)
				} else {
					c.JSON(400, err.Error())
				}
			})

		} else {
			c.JSON(400, err.Error())
		}

	})

	r.NoRoute(func(c *gin.Context) {
		c.File(fmt.Sprintf("%s/index.html", _server_root))
	})

	log.WithFields(log.Fields{
		"Port":    port,
		"Webroot": server_root,
	}).Debug("Interface")
	r.Run(fmt.Sprintf(":%d", port))

}
