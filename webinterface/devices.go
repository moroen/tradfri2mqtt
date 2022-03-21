package webinterface

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/moroen/tradfri2mqtt/tradfri"
	"github.com/sirupsen/logrus"
)

type DevicePayload struct {
	State bool `json:"state"`
}

func DeviceRoutes(r *gin.Engine) error {
	r.GET("/devices", func(c *gin.Context) {
		dev := tradfri.GetDevicesList()
		c.JSON(http.StatusOK, dev)
	})

	r.GET("/devices/:deviceid", func(c *gin.Context) {
		if deviceid, err := strconv.Atoi(c.Param("deviceid")); err == nil {
			tradfri.GetDevice(deviceid, func(dev *tradfri.TradfriDevice, err error) {
				c.JSON(http.StatusOK, dev)
			})
		} else {
			c.JSON(http.StatusBadRequest, PostResponse{Status: "Error", Error: err.Error()})
		}
	})

	r.POST("/api/devices/:deviceid/set/state", func(c *gin.Context) {
		var cmd DevicePayload

		stat := make(chan gin.H)

		if deviceid, err := strconv.Atoi(c.Param("deviceid")); err == nil {
			if c.ShouldBind(&cmd); err == nil {

				go tradfri.GetDevice(deviceid, func(d *tradfri.TradfriDevice, err error) {
					if err == nil {
						d.SetState(cmd.State, func(msg []byte, err error) {
							fmt.Println("SetState Done")
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
						fmt.Println("Fails here")
						logrus.WithFields(logrus.Fields{
							"Error": err.Error(),
						}).Error("API - SetState")
						stat <- gin.H{"Status": "Error", "Error": "Unable to find device"}
					}
				})

			} else {
				fmt.Println(err.Error())
				c.JSON(http.StatusOK, gin.H{"Status": "Error", "Error": err.Error()})
				return
			}
		}
		c.JSON(http.StatusOK, <-stat)
	})

	return nil
}
