package tradfri

import (
	"fmt"
	"sync"

	coap "github.com/moroen/go-tradfricoap"
	"github.com/moroen/tradfri2mqtt/mqttclient"
	"github.com/moroen/tradfri2mqtt/settings"
	log "github.com/sirupsen/logrus"
)

func Start(wg *sync.WaitGroup, status_channel chan (error)) {
	conf := settings.GetCoapConfig(true)

	log.Debug("Connecting to gateway")
	if err := coap.ConnectGateway(conf); err != nil {
		fmt.Println("Shite")
	} else {
		log.Info(fmt.Sprintf("Connected to tradfri at [tcp://%s:%d]", conf.Gateway, 5684))
		if _, err := coap.Observe(wg, mqttclient.SendState, conf.KeepAlive); err != nil {
			log.Error(err.Error())
		}
	}
}

func Stop() {
	coap.ObserveStop()
}
