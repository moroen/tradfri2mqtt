package mqttclient

import (
	"encoding/json"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	coap "github.com/moroen/go-tradfricoap"
	"github.com/moroen/tradfri2mqtt/settings"
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

	var gwconf MQTTgwConfig
	err := json.Unmarshal(msg.Payload(), &gwconf)
	if err != nil {
		log.Fatalln(err.Error())
	}

	switch s[2] {
	case "gwconfig":
		gwConfig(gwconf, client)
	}

}

func gwConfig(gwconf MQTTgwConfig, client mqtt.Client) error {
	cfg := settings.GetConfig(false)

	if psk, err := coap.GetNewPSK(gwconf.Ip, gwconf.Key); err == nil {
		cfg.Tradfri.Gateway = gwconf.Ip
		cfg.Tradfri.Identity = psk.Ident
		cfg.Tradfri.Passkey = psk.Key
		fmt.Println("Writing config")
		err = settings.WriteConfig(&cfg)
		if err != nil {
			log.Error(err.Error())
			if statusJson, err := json.Marshal(MQTTstatus{Status: "Error", Error: err.Error()}); err == nil {
				return SendTopic("tradfri/cmd/status/gwconfig", statusJson)
			}
		}

		if statusJson, err := json.Marshal(MQTTstatus{Status: "Ok"}); err == nil {
			log.Debug("gwConfig ok")
			return SendTopic("tradfri/cmd/status/gwconfig", statusJson)
		}
		return nil
	} else {
		log.Fatal(err.Error())
		return err
	}
}
