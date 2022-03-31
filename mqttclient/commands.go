package mqttclient

import (
	"encoding/json"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/moroen/tradfri2mqtt/tradfri"
)

type MQTTstatus struct {
	Status string
	Error  string
}

func Command(client mqtt.Client, msg mqtt.Message) {
	log.WithFields(log.Fields{
		"topic":   msg.Topic(),
		"payload": string(msg.Payload()),
	}).Debug("Received command message")

	s := strings.Split(msg.Topic(), "/")

	switch s[2] {
	case "discover":
		tradfri.Discover(true)
	case "test":
		tradfri.Test()
	case "reboot":
		tradfri.RebootGateway()
	default:
		str := map[string]interface{}{"Unknown command": s[2]}

		if payload, err := json.Marshal(str); err == nil {
			SendTopic(fmt.Sprintf("%s/status/error", viper.GetString("mqtt.discoverytopic")), payload, false)
		} else {
			log.WithFields(log.Fields{
				"Error": err.Error(),
			}).Error("MQTT - Command")
		}
	}

}
