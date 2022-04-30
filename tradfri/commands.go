package tradfri

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type MQTTStateMessage struct {
	Color struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	} `json:"color"`
	State      interface{} `json:"state"`
	Brightness int         `json:"brightness"`
	ColorTemp  int         `json:"color_temp"`
	Position   int         `json:"position"`
}

type MQTTpayload struct {
	Value int64
}

type MQTTStatePayload struct {
	State string
}

type MQTTboolPayload struct {
	Value string
}

func ParseMessage(topic string, payload []byte) (deviceid int64, state int, level int, x int64, y int64, color_temp int, err error) {

	message := MQTTStateMessage{Brightness: -1, Position: -1}
	message.ColorTemp = -1
	message.Color.X = -1
	message.Color.Y = -1

	if err = json.Unmarshal(payload, &message); err == nil {
		s := strings.Split(topic, "/")

		deviceid, err = strconv.ParseInt(s[1], 10, 64)
		if err != nil {
			log.Fatalln(err.Error())
			return 0, 0, 0, -1, -1, -1, err
		}

		switch message.State.(type) {
		case string:
			switch strings.ToUpper(fmt.Sprintf("%v", message.State)) {
			case "ON":
				state = 1
			case "OFF":
				state = 0
			default:
				state = -1
			}
		case bool:
			if message.State.(bool) {
				state = 1
			} else {
				state = 0
			}
		default:
			state = -1
		}

		if message.Color.X != -1 {
			x = int64(message.Color.X * 65535)
		} else {
			x = -1
		}

		if message.Color.Y != -1 {
			y = int64(message.Color.Y * 65535)
		} else {
			y = -1
		}

		if message.Position != -1 {
			return deviceid, state, message.Position, x, y, message.ColorTemp, nil
		} else {
			return deviceid, state, message.Brightness, x, y, message.ColorTemp, nil
		}
	} else {
		log.WithFields(log.Fields{
			"Error": err.Error(),
		}).Error("PrseMessage - Unmarshal")
		return 0, 0, 0, -1, -1, -1, err
	}
}

func State(topic string, payload []byte) {
	state := -1

	log.WithFields(log.Fields{
		"topic":   topic,
		"payload": string(payload),
	}).Debug("Received state message")

	if bValue, err := strconv.ParseBool(string(payload)); err == nil {
		if bValue {
			state = 1
		} else {
			state = 0
		}
	} else {
		var payloadStruct MQTTboolPayload
		if err := json.Unmarshal(payload, &payloadStruct); err == nil {
			if payloadStruct.Value == "true" {
				state = 1
			} else {
				state = 0
			}
		}
	}

	s := strings.Split(topic, "/")

	deviceid, err := strconv.ParseInt(s[1], 10, 64)
	if err != nil {
		log.Fatalln(err.Error())
	}

	GetDevice(int(deviceid), false, func(d *TradfriDevice, err error) {
		d.SetState(state != 0, func(msg []byte, err error) {
			if err != nil {
				log.WithFields(log.Fields{
					"Error": err.Error(),
				}).Error("MQTT - Set State")
			}
		})
	})

}

func Blind(topic string, payload []byte) {
	log.WithFields(log.Fields{
		"topic":   topic,
		"payload": string(payload),
	}).Debug("Received blind message")

	deviceid, _, value, _, _, _, err := ParseMessage(topic, payload)
	if err != nil {
		log.Error(err.Error())
		return
	}

	// value = 100 - value

	log.WithFields(log.Fields{
		"DeviceID": deviceid,
		"Position": value,
	}).Debug("Handlers - Blind")

	GetDevice(int(deviceid), false, func(d *TradfriDevice, err error) {
		if err == nil {
			d.SetBlind(value, func(msg []byte, err error) {
				if err != nil {
					log.WithFields(log.Fields{
						"Error": err.Error(),
					}).Error("MQTT - Blind")
				}
			})
		} else {
			log.WithFields(log.Fields{
				"Error": err.Error(),
			}).Error("MQTT - Blind - Unable to get device")
		}
	})

}

func Dimmer(topic string, payload []byte) {
	// fmt.Printf("Received dimmer message: %s from topic: %s\n", msg.Payload(), msg.Topic())

	log.WithFields(log.Fields{
		"topic":   topic,
		"payload": string(payload),
	}).Debug("Received dimmer message")

	deviceid, state, level, x, y, col_temp, err := ParseMessage(topic, payload)
	if err != nil {
		log.Error(err.Error())
		return
	}

	log.WithFields(log.Fields{
		"State":    state,
		"Level":    level,
		"X":        x,
		"Y":        y,
		"Col_Temp": col_temp,
	}).Debug("MQTT-handlers - Dimmer")

	GetDevice(int(deviceid), false, func(d *TradfriDevice, err error) {

		if state != -1 {
			d.SetState(state != 0, func(msg []byte, err error) {
				if err != nil {
					log.WithFields(log.Fields{
						"Error": err.Error(),
					}).Error("MQTT - Dimmer - SetState")
				}
			})

		}

		if level != -1 {
			if d.LightControl != nil {
				d.SetLevel(level, func(msg []byte, err error) {
					if err != nil {
						log.WithFields(log.Fields{
							"Error": err.Error(),
						}).Error("MQTT - Set State")
					}
				})
			} else if d.BlindControl != nil {
				d.SetBlind(level, func(msg []byte, err error) {
					if err != nil {
						log.WithFields(log.Fields{
							"Error": err.Error(),
						}).Error("MQTT - Set Blind")
					}
				})

			}
		}

		if x != -1 {
			d.SetXY(int(x), int(y), func(msg []byte, err error) {
				if err != nil {
					log.WithFields(log.Fields{
						"Error": err.Error(),
					}).Error("MQTT - Set XY")
				}
			})
		}

		if col_temp != -1 {
			log.WithFields(log.Fields{
				"col_temp": col_temp,
			}).Debug("MQTT-handlers - Dimmer - SetColorTemp")

			var hex string

			if col_temp >= 350 {
				hex = "efd275"
			} else if col_temp < 280 {
				hex = "f5faf6"
			} else {
				hex = "f1e0b5"
			}

			d.SetHex(hex, func(msg []byte, err error) {
				if err != nil {
					log.WithFields(log.Fields{
						"Error": err.Error(),
					}).Error("MQTT - Set HEX")
				}
			})
		}

		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Error("Dimmer message error")
		}
	})
}

func RebootGateway() {
	uri := fmt.Sprintf("%s/%s", attrGateway, attrReboot)
	ctx, done := context.WithTimeout(context.Background(), 2*time.Second)
	defer done()

	log.Debug("Rebooting gateway")

	if _connection.ObserveDone != nil {
		_connection.ObserveDone()
		_connection.ObserveWaitGroup.Wait()
	}

	_connection.POST(ctx, uri, "1", func(msg []byte, err error) {
		log.WithFields(log.Fields{
			"Message": string(msg),
		}).Debug("commands.RebootGateway")
	})
	_connection.Disconnect()
}
