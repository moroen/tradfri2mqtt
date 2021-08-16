package main

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	coap "github.com/moroen/go-tradfricoap"
	"github.com/moroen/tradfri2mqtt/settings"
	"github.com/moroen/tradfri2mqtt/tradfri"
)

var status_channel chan (error)

func main() {

	log.SetLevel(log.DebugLevel)
	coap.SetCoapRetry(2, 1)

	latest_restart := time.Now()

	status_channel = make(chan error)

	conf := settings.GetConfig(false)
	settings.WriteConfig(&conf)

	// go mqttclient.Start(status_channel)
	// go tradfri.Start(status_channel)

	// time.Sleep(2 * time.Second)
	//coap.ObserveRestart(true)
	/*
		select {
		case err := <-status_channel:
			fmt.Println("Error")
			log.Error(err.Error())
		}
	*/

	var err error

	go Interface_Server()

	for err == nil {
		select {
		case err = <-status_channel:
			fmt.Println(err.Error())
			fmt.Println("Done")
		default:
			if conf.Tradfri.KeepAlive != 0 {
				if diff := time.Since(latest_restart); diff > (time.Second * time.Duration(conf.Tradfri.KeepAlive)) {
					go tradfri.ReStart()
					latest_restart = time.Now()
				}
			}
		}
	}
}
