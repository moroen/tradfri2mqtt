package tradfri

import (
	"context"
	"errors"
	"fmt"
	"sync"

	coap "github.com/moroen/go-tradfricoap"
	"github.com/moroen/tradfri2mqtt/settings"
	log "github.com/sirupsen/logrus"
)

var ErrorTradfriConnectionRefused = errors.New("Tradfri Error: Connection refused")

var _wg *sync.WaitGroup
var _ctx context.Context
var _stop func()

var cfg coap.GatewayConfig

func onConnect() {
	log.Info(fmt.Sprintf("Tradfri: !!! Connected to gateway at [tcp://%s:%s]", cfg.Gateway, "5684"))
}

func Start(wg *sync.WaitGroup, status_channel chan (error), connectDelay int) {
	log.Info("Tradfri: Starting")
	_wg = wg
	_wg.Add(1)

	_ctx, _stop = context.WithCancel(context.Background())

	cfg = settings.GetCoapConfig(false)

	coap.ConnectGateway(_ctx, cfg, onConnect,
		nil,
		nil,
		nil)
	/*
		conf := settings.GetCoapConfig(false)
		coap.SetConfig(conf)

		// ctx, stop := context.WithCancel(context.Background())
		// _stop = stop
		coap.Observe(mqttclient.SendState, conf.KeepAlive)
	*/
}

func Stop() {
	log.Info("Tradfri: Stopping")
	coap.ObserveStop()
	_wg.Done()
	log.Info("Tradfri: Stoppped")
}
