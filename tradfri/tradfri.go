package tradfri

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	coap "github.com/moroen/go-tradfricoap"
	"github.com/moroen/tradfri2mqtt/mqttclient"
	"github.com/moroen/tradfri2mqtt/settings"
	log "github.com/sirupsen/logrus"
)

var ErrorTradfriConnectionRefused = errors.New("Tradfri Error: Connection refused")

var _wg *sync.WaitGroup

var _stop func()

func connect_tradfri(ctx context.Context, conf coap.GatewayConfig) error {
	ticker := time.NewTicker(5 * time.Second)
	for {
		if err := coap.ConnectGateway(conf); err == nil {
			log.Info(fmt.Sprintf("Connected to tradfri at %s:5684", conf.Gateway))
			coap.Observe(mqttclient.SendState, conf.KeepAlive)
			return nil
		} else {
			log.Error("Tradfri: Unable to connect to gateway")
		}
		select {
		case <-ticker.C:
			break
		case <-ctx.Done():
			fmt.Println("tock")
			return errors.New("tradfri: tradfri connection canceled")
		}
	}
}

func Start(wg *sync.WaitGroup, status_channel chan (error), connectDelay int) {
	log.Info("Tradfri: Starting")
	_wg = wg

	_wg.Add(1)

	conf := settings.GetCoapConfig(false)
	coap.SetConfig(conf)

	// ctx, stop := context.WithCancel(context.Background())
	// _stop = stop
	coap.Observe(mqttclient.SendState, conf.KeepAlive)
}

func Stop() {
	log.Info("Tradfri: Stopping")
	coap.ObserveStop()
	_wg.Done()
	log.Info("Tradfri: Stoppped")
}
