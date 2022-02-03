package webinterface

import (
	"fmt"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/moroen/tradfri2mqtt/tradfri"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

var status_channel chan (error)

var logLevelMap = map[string]logrus.Level{
	"trace": logrus.TraceLevel,
	"debug": logrus.DebugLevel,
	"info":  logrus.InfoLevel,
	"warn":  logrus.WarnLevel,
	"error": logrus.ErrorLevel,
}

var r *gin.Engine

type arguments struct {
	LogLevel       string
	BindAddress    string
	BindPort       int
	StaticContents string
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

func Interface_Server(server_root string, status_channel chan (error)) {

	status_channel = status_channel
	_server_root = server_root
	fmt.Print(server_root)

	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&logrus.JSONFormatter{})

	r = gin.Default()
	r.Use(CORS())

	r.Use(static.Serve("/", static.LocalFile(_server_root, false)))

	r.GET("/api/v1/hello", func(c *gin.Context) {
		c.JSON(200, `{"message":"hello, hello, hello"}`)
	})

	r.GET("/api/settings", func(c *gin.Context) {
		// conf := settings.GetConfig(false)

		c.JSON(200, "")
	})

	r.POST("/api/settings", func(c *gin.Context) {
		// conf := settings.GetConfig(false)

		fmt.Println(c.Request.Body)

		/*
			if err := c.ShouldBind(&conf); err == nil {
				ret := conf
				ret.Save()
				c.JSON(200, ret)
				err := settings.ErrConfigIsDirty
				status_channel <- err
			} else {
				c.JSON(400, err.Error())
			}
		*/
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

	/*
		r.NoRoute(func(c *gin.Context) {
			c.File("./www/index.html")
		})
	*/

	r.Run(":8321")
}
