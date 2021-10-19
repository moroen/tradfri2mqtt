package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	coap "github.com/moroen/go-tradfricoap"
	"github.com/moroen/tradfri2mqtt/mqttclient"
	"github.com/moroen/tradfri2mqtt/settings"
	"github.com/moroen/tradfri2mqtt/tradfri"
)

var status_channel chan (error)

func main() {

	var wg sync.WaitGroup

	log.SetLevel(log.DebugLevel)
	// log.SetReportCaller(true)
	coap.SetCoapRetry(2, 1)

	var err error
	var tradfri_control chan (error)

	// latest_restart := time.Now()

	status_channel = make(chan error)

	settings.GetConfig(false)

	go mqttclient.Start(&wg, status_channel)

	go tradfri.Start(&wg, status_channel)
	// go Interface_Server(conf.Interface.ServerRoot)
	// go mqttclient.Do_Test(&wg)
	// time.Sleep(2 * time.Second)
	//coap.ObserveRestart(true)
	/*
		select {
		case err := <-status_channel:
			fmt.Println("Error")
			log.Error(err.Error())
		}
	*/

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		for err == nil {
			select {
			case <-c:
				fmt.Println("Sig catched")
				tradfri.Stop()
				mqttclient.Stop()
				wg.Wait()
				// <-tradfri_controll
				os.Exit(1)
			case err = <-tradfri_control:
				fmt.Println(err.Error())
				/*
					case err = <-status_channel:
						log.Debug(err.Error())
						if err == coap.ErrorNoConfig {
							err = nil
						} else if err == errors.ErrorConfigStale {
							conf = settings.GetConfig(true)
							mqttclient.Restart()
							// tradfri.ReStart()
							err = nil
						} else {
							fmt.Println("Done")
						}
				*/
			}
		}
	}()
	time.Sleep(5 * time.Second)
	// coap.ObserveRestart(false)
	wg.Wait()
}
