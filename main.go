package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	coap "github.com/moroen/go-tradfricoap"
	"github.com/moroen/tradfri2mqtt/mqttclient"
	"github.com/moroen/tradfri2mqtt/tradfri"
)

var status_channel chan (error)

func main() {

	log.SetLevel(log.DebugLevel)
	coap.SetCoapRetry(2, 1)

	status_channel = make(chan error)

	// conf := settings.GetConfig(false)

	go mqttclient.Start(status_channel)
	go tradfri.Start(status_channel)

	// time.Sleep(2 * time.Second)
	//coap.ObserveRestart(true)
	/*
		select {
		case err := <-status_channel:
			fmt.Println("Error")
			log.Error(err.Error())
		}
	*/

	err := <-status_channel
	fmt.Println(err.Error())
	fmt.Println("Done")
}
