package tradfri

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/buger/jsonparser"
	coap "github.com/moroen/gocoap/v5"
	"github.com/moroen/tradfri2mqtt/settings"
	log "github.com/sirupsen/logrus"
)

var ErrorTradfriConnectionRefused = errors.New("Tradfri Error: Connection refused")

var _wg *sync.WaitGroup

// Observe
var _wgObserve sync.WaitGroup
var _ctxObserve context.Context
var _stopObserve func()

var cfg settings.Config

var MQTTSendTopic func(string, []byte, bool) error

var _connection coap.CoapDTLSConnection

var _mqtt_command_topic string
var _mqtt_discovery_topic string

func onConnect() {
	log.Info(fmt.Sprintf("Tradfri: Connected to gateway at [tcp://%s:%s]", cfg.Tradfri.Gateway, "5684"))
	MQTTSendTopic("tradfri/status", []byte("Connected"), false)
	Discover(false)
	cfg = settings.GetConfig(false)
	Observe()
}

func SetConnecion(conn coap.CoapDTLSConnection) {
	_connection = conn
}

func Start(wg *sync.WaitGroup, status_channel chan (error)) {
	log.Info("Tradfri: Starting")
	_wg = wg
	wg.Add(1)
	_devices.Init()

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
		UseQueue:  true,
		KeepAlive: cfg.Tradfri.KeepAlive,
	}
	go _connection.Connect()
}

func Stop() {
	log.Info("Tradfri: Stopping")
	if _stopObserve != nil {
		_stopObserve()
	}
	_wgObserve.Wait()
	if err := _connection.Disconnect(); err != nil {
		log.WithFields(log.Fields{
			"Error": err.Error(),
		}).Error("Tradfri - Stop")
	}

	log.Info("Tradfri: Stopped")
	_wg.Done()
}

func Test() {

	_devices.GetDeviceInfo(65549, func(device *TradfriDevice, err error) {
		if uri, payload, err := device.SetState(0); err == nil {
			ctx, done := context.WithTimeout(context.Background(), 2*time.Second)
			defer done()
			_connection.PUT(ctx, uri, payload, func(msg []byte, err error) {
				fmt.Println(string(msg))
			})
		} else {
			fmt.Println(err.Error())
		}
	})
	/*
		ctx, done := context.WithTimeout(context.Background(), time.Second*2)
		defer done()

		_connection.GET(ctx, "15001/65560", func(msg []byte, err error) {

			fmt.Println(string(msg))

			var device TradfriDevice
			if err := json.Unmarshal(msg, &device); err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Printf("%+v\n", device)
				fmt.Println(device.Id)
				fmt.Println(device.LightControl[0].State)
			}
		})
	*/
}

func Observe() {

	_ctxObserve, _stopObserve = context.WithCancel(context.Background())

	ctx, done := context.WithTimeout(context.Background(), 3*time.Second)
	defer done()

	_connection.GET(ctx, uriDevices, func(msg []byte, err error) {
		if _, err = jsonparser.ArrayEach(msg, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			if res, err := jsonparser.GetInt(value); err == nil {
				go _connection.Observe(fmt.Sprintf("%s/%d", uriDevices, res), func(msg []byte, err error) {
					// go _connection.Observe(_ctxObserve, &_wgObserve, fmt.Sprintf("%s/%d", uriDevices, res), func(msg []byte, err error) {
					if err != nil {
						log.Error("Test - Observe failed")
					} else {
						SendState(msg)
					}
				})
			}
		}); err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Error("Tradfri - Observe")
		}
	})

}
