package tradfri

import (
	"context"
	"errors"
	"fmt"
	"sync"

	coap "github.com/moroen/gocoap/v5"
	"github.com/moroen/tradfri2mqtt/settings"
	log "github.com/sirupsen/logrus"
)

var ErrorTradfriConnectionRefused = errors.New("Tradfri Error: Connection refused")

var _wg *sync.WaitGroup
var _ctx context.Context
var _stop func()

var cfg settings.Config

var MQTTSendTopic func(string, []byte, bool) error

var _connection coap.CoapDTLSConnection

var _mqtt_command_topic string
var _mqtt_discovery_topic string

func onConnect() {
	log.Info(fmt.Sprintf("Tradfri: Connected to gateway at [tcp://%s:%s]", cfg.Tradfri.Gateway, "5684"))
	MQTTSendTopic("tradfri/status", []byte("Connected"), false)
	Discover(false)
}

func Start(wg *sync.WaitGroup, status_channel chan (error)) {
	log.Info("Tradfri: Starting")
	_wg = wg
	_wg.Add(1)

	_devices.Init()

	_ctx, _stop = context.WithCancel(context.Background())

	cfg = settings.GetConfig(false)
	_mqtt_command_topic = cfg.Mqtt.CommandTopic
	_mqtt_discovery_topic = cfg.Mqtt.DiscoveryTopic

	_connection = coap.CoapDTLSConnection{
		Host:      cfg.Tradfri.Gateway,
		Port:      5684,
		Ident:     cfg.Tradfri.Identity,
		Key:       cfg.Tradfri.Passkey,
		OnConnect: onConnect,
		OnConnectionFailed: func() {
			log.Info(fmt.Sprintf("Tradfri: Unable to connected to gateway at [tcp://%s:%s]", cfg.Tradfri.Gateway, "5684"))
		},
		OnDisconnect: func() {
			log.Info(fmt.Sprintf("Tradfri: Disconnected from gateway at [tcp://%s:%s]", cfg.Tradfri.Gateway, "5684"))
		},
		UseQueue: true,
	}
	_connection.Connect()
}

func Stop() {
	log.Info("Tradfri: Stopping")
	// coap.ObserveStop()
	_connection.Disconnect()
	_wg.Done()
	log.Info("Tradfri: Stopped")
}

func Test() {
	_devices.GetDeviceInfo(65552, func(msg []byte, err error) {
		fmt.Println(string(msg))
	})
}
