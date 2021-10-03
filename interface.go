package main

import (
	"fmt"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/moroen/tradfri2mqtt/settings"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

var logLevelMap = map[string]logrus.Level{
	"trace": logrus.TraceLevel,
	"debug": logrus.DebugLevel,
	"info":  logrus.InfoLevel,
	"warn":  logrus.WarnLevel,
	"error": logrus.ErrorLevel,
}

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
	Ident string `json:"ident"`
	Key   string `json:"key"`
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

func Interface_Server(server_root string) {

	_server_root = server_root
	fmt.Print(server_root)

	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&logrus.JSONFormatter{})

	r := gin.Default()
	r.Use(CORS())

	r.Use(static.Serve("/", static.LocalFile(_server_root, false)))

	r.GET("/api/v1/hello", func(c *gin.Context) {
		c.JSON(200, `{"message":"hello, hello, hello"}`)
	})

	r.GET("/api/settings", func(c *gin.Context) {
		conf := settings.GetConfig(false)

		c.JSON(200, conf)
	})

	r.POST("/api/settings", func(c *gin.Context) {
		var conf settings.Config

		if err := c.ShouldBind(&conf); err == nil {
			fmt.Println(conf)
			ret := conf
			c.JSON(200, ret)
		} else {
			c.JSON(400, err.Error())
		}
	})

	r.POST("/api/getPSK", func(c *gin.Context) {
		conf := settings.GetConfig(false)
		var data TradfriFields

		if err := c.ShouldBind(&data); err == nil {
			fmt.Println(data)
			ret := TradfriPSk{Ident: conf.Tradfri.Identity, Key: conf.Tradfri.Passkey}
			c.JSON(200, ret)
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
