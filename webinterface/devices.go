package webinterface

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/moroen/tradfri2mqtt/tradfri"
	"github.com/sirupsen/logrus"
)

type DevicePayload struct {
	Command string `json:"command"`
	State   bool   `json:"state"`
	Level   int    `json:"level"`
	Hex     string `json:"hex"`
}

func DeviceRoutes(r *gin.Engine) error {
	r.GET("api/devices", func(c *gin.Context) {
		dev := tradfri.GetDevicesList()
		c.JSON(http.StatusOK, dev)
	})

	r.GET("api/devices/:deviceid", func(c *gin.Context) {
		if deviceid, err := strconv.Atoi(c.Param("deviceid")); err == nil {
			tradfri.GetDevice(deviceid, func(dev *tradfri.TradfriDevice, err error) {
				c.JSON(http.StatusOK, dev)
			})
		} else {
			c.JSON(http.StatusBadRequest, PostResponse{Status: "Error", Error: err.Error()})
		}
	})

	r.POST("/api/devices/:deviceid/set", func(c *gin.Context) {
		var cmd DevicePayload

		stat := make(chan gin.H)

		if deviceid, err := strconv.Atoi(c.Param("deviceid")); err == nil {
			if c.ShouldBind(&cmd); err == nil {
				switch cmd.Command {
				case "state":
					go setState(deviceid, cmd, stat)
				case "level":
					go setLevel(deviceid, cmd, stat)
				default:
					c.JSON(http.StatusOK, gin.H{"Status": "Error", "Error": "Unknown command"})
					return
				}

			} else {
				logrus.WithFields(logrus.Fields{
					"Error": err.Error(),
				}).Error("devices.route.set failed")
				c.JSON(http.StatusOK, gin.H{"Status": "Error", "Error": err.Error()})
				return
			}
		}
		c.JSON(http.StatusOK, <-stat)
	})

	return nil
}

func setLevel(deviceid int, cmd DevicePayload, stat chan gin.H) {
	tradfri.GetDevice(deviceid, func(d *tradfri.TradfriDevice, err error) {
		if err == nil {
			d.SetLevel(cmd.Level, func(msg []byte, err error) {
				if err == nil {
					stat <- gin.H{"Status": "Ok"}
				} else {
					stat <- gin.H{"Status": "Error", "Error": err.Error()}
					logrus.WithFields(logrus.Fields{
						"Error": err.Error(),
					}).Error("API - SetState")
				}
			})
		} else {
			logrus.WithFields(logrus.Fields{
				"Error": err.Error(),
			}).Error("API - SetState")
			stat <- gin.H{"Status": "Error", "Error": "Unable to find device"}
		}
	})
}

func setState(deviceid int, cmd DevicePayload, stat chan gin.H) {
	tradfri.GetDevice(deviceid, func(d *tradfri.TradfriDevice, err error) {
		if err == nil {
			d.SetState(cmd.State, func(msg []byte, err error) {
				if err == nil {
					stat <- gin.H{"Status": "Ok"}
				} else {
					stat <- gin.H{"Status": "Error", "Error": err.Error()}
					logrus.WithFields(logrus.Fields{
						"Error": err.Error(),
					}).Error("API - SetState")
				}
			})
		} else {
			logrus.WithFields(logrus.Fields{
				"Error": err.Error(),
			}).Error("API - SetState")
			stat <- gin.H{"Status": "Error", "Error": "Unable to find device"}
		}
	})
}
