package mqttclient

import (
	"fmt"

	"github.com/moroen/tradfri2mqtt/tradfri"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTgwConfig struct {
	Ip  string
	Key string
}

type QueueItem struct {
	Client  mqtt.Client
	Message mqtt.Message
	Handler func(mqtt.Client, mqtt.Message)
	Retries int
}

var ItemQueue []QueueItem

func Subscribe(client mqtt.Client, status_channel chan (error)) {
	_status_channel = status_channel

	if client == nil {
		log.Fatal("Unable to subscribe, client is nil")
		return
	}

	if token := client.Subscribe(fmt.Sprintf("%s/+/WS/set", viper.GetString("mqtt.commandtopic")), 0, SetHex); token.Wait() && token.Error() != nil {
		log.Print(token.Error())
	}

	if token := client.Subscribe(fmt.Sprintf("%s/+/CWS/set", viper.GetString("mqtt.commandtopic")), 0, SetHex); token.Wait() && token.Error() != nil {
		log.Print(token.Error())
	}

	if token := client.Subscribe(fmt.Sprintf("%s/+/dimmer/set", viper.GetString("mqtt.commandtopic")), 0, Dimmer); token.Wait() && token.Error() != nil {
		log.Print(token.Error())
	}

	if token := client.Subscribe(fmt.Sprintf("%s/+/blind/set", viper.GetString("mqtt.commandtopic")), 0, Blind); token.Wait() && token.Error() != nil {
		log.Print(token.Error())
	}

	if token := client.Subscribe(fmt.Sprintf("%s/+/switch/set", viper.GetString("mqtt.commandtopic")), 0, State); token.Wait() && token.Error() != nil {
		log.Print(token.Error())
	}

	if token := client.Subscribe(fmt.Sprintf("%s/cmd/#", viper.GetString("mqtt.discoverytopic")), 0, Command); token.Wait() && token.Error() != nil {
		log.Print(token.Error())
	}
}

func State(client mqtt.Client, msg mqtt.Message) {
	tradfri.State(msg.Topic(), msg.Payload())
}

func SetHex(client mqtt.Client, msg mqtt.Message) {
}

func Blind(client mqtt.Client, msg mqtt.Message) {
	tradfri.Blind(msg.Topic(), msg.Payload())
}

func Dimmer(client mqtt.Client, msg mqtt.Message) {
	tradfri.Dimmer(msg.Topic(), msg.Payload())
}
