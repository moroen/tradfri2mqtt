package tradfri

import (
	"sync"

	coap "github.com/moroen/go-tradfricoap"
	"github.com/moroen/tradfri2mqtt/mqttclient"
	"github.com/moroen/tradfri2mqtt/settings"
	log "github.com/sirupsen/logrus"
)

func Start(wg *sync.WaitGroup, status_channel chan (error)) {
	conf := settings.GetCoapConfig(true)
	coap.SetConfig(conf)

	_, err := coap.Observe(wg, mqttclient.SendState, conf.KeepAlive)
	if err != nil {
		log.Error(err.Error())
	}
}

func Stop() {
	coap.ObserveStop()
}
