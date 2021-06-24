package tradfri

import (
	coap "github.com/moroen/go-tradfricoap"
	"github.com/moroen/tradfri2mqtt/settings"
	log "github.com/sirupsen/logrus"
)

var status_channel chan (error)

func SetStatusChannel(c chan (error)) {
	status_channel = c
}

func Start(reloadConfig bool) error {
	coapConfig := settings.GetCoapConfig(reloadConfig)
	coap.SetConfig(coapConfig)
	// go coap.Observe(messages.Show, status_channel)
	return nil
}

func ReStart() error {
	log.Debug("Tradfri Restart")
	coap.ObserveStop()
	coap.CloseDTLSConnection()
	Start(true)
	return nil
}
