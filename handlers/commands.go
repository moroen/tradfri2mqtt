package handlers

import (
	"encoding/json"
	"strings"

	log "github.com/sirupsen/logrus"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	coap "github.com/moroen/go-tradfricoap"
	"github.com/moroen/tradfri2mqtt/settings"
)

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
		gwConfig(gwconf)
	}

}

func gwConfig(gwconf MQTTgwConfig) error {
	cfg := settings.GetConfig()

	if psk, err := coap.GetNewPSK(gwconf.Ip, gwconf.Key); err == nil {
		cfg.Tradfri.Gateway = gwconf.Ip
		cfg.Tradfri.Identity = psk.Ident
		cfg.Tradfri.Passkey = psk.Key
		settings.WriteConfig(&cfg)
		return nil
	} else {
		return err
	}
}
