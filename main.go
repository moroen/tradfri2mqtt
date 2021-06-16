package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	coap "github.com/moroen/go-tradfricoap"
	"github.com/moroen/tradfri2mqtt/handlers"
	"github.com/moroen/tradfri2mqtt/messages"
	"github.com/moroen/tradfri2mqtt/settings"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	opts := client.OptionsReader()
	log.Info(fmt.Sprintf("Connected to broker at %s", opts.Servers()))
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Info(fmt.Sprintf("Connect lost: %v", err))
}

func main() {

	log.SetLevel(log.DebugLevel)
	coap.SetCoapRetry(2, 1)

	c := make(chan error)

	conf := settings.GetConfig()
	coapConfig := settings.GetCoapConfig()
	coap.SetConfig(coapConfig)

	var broker = conf.Mqtt.Host
	var port = 1883
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID("go_mqtt_client")
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	_client := mqtt.NewClient(opts)
	if token := _client.Connect(); token.Wait() && token.Error() != nil {
		log.Error(token.Error())
		return
	}

	messages.SetClientConnection(_client)
	handlers.Subscribe(_client)
	go coap.Observe(messages.Show, c)
	// time.Sleep(2 * time.Second)
	//coap.ObserveRestart(true)
	select {
	case err := <-c:
		log.Error(err.Error())
	}
	<-c

	_client.Disconnect(250)
}
