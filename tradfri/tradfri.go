package tradfri

import (
	coap "github.com/moroen/go-tradfricoap"
	"github.com/moroen/tradfri2mqtt/mqttclient"
	"github.com/moroen/tradfri2mqtt/settings"
	log "github.com/sirupsen/logrus"
)

var _status_channel chan (error)

func Start(status_channel chan (error)) error {
	_status_channel = status_channel
	coapConfig := settings.GetCoapConfig(false)
	coap.SetConfig(coapConfig)
	go coap.Observe(mqttclient.Show, status_channel)
	return nil
}

func ReStart() error {
	log.Debug("Tradfri Restart")
	coap.ObserveStop()
	coap.CloseDTLSConnection()
	Start(_status_channel)
	return nil
}
