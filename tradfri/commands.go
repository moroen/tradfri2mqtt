package tradfri

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
)

func State(deviceID int, state int) {

	_devices.GetDeviceInfo(deviceID, func(device *TradfriDevice, err error) {
		if uri, payload, err := device.SetState(state); err == nil {
			ctx, done := context.WithTimeout(context.Background(), 2*time.Second)
			defer done()
			_connection.PUT(ctx, uri, payload, func(msg []byte, err error) {})
		} else {
			log.WithFields(log.Fields{
				"Error": err.Error(),
			}).Error("Tradfri - State")
		}
	})

}

func Level(deviceID int, state int) {

	_devices.GetDeviceInfo(deviceID, func(device *TradfriDevice, err error) {
		if uri, payload, err := device.SetLevel(state); err == nil {
			ctx, done := context.WithTimeout(context.Background(), 2*time.Second)
			defer done()
			_connection.PUT(ctx, uri, payload, func(msg []byte, err error) {})
		} else {
			log.WithFields(log.Fields{
				"Error": err.Error(),
			}).Error("Tradfri - State")
		}
	})

}

func SetXY(deviceID int, x int, y int) {

	_devices.GetDeviceInfo(deviceID, func(device *TradfriDevice, err error) {
		if uri, payload, err := device.SetXY(x, y); err == nil {
			ctx, done := context.WithTimeout(context.Background(), 2*time.Second)
			defer done()
			_connection.PUT(ctx, uri, payload, func(msg []byte, err error) {})
		} else {
			log.WithFields(log.Fields{
				"Error": err.Error(),
			}).Error("Tradfri - SetXY")
		}
	})

}

func SetHex(deviceID int, hex string) {

	_devices.GetDeviceInfo(deviceID, func(device *TradfriDevice, err error) {
		if uri, payload, err := device.SetHex(hex); err == nil {
			ctx, done := context.WithTimeout(context.Background(), 2*time.Second)
			defer done()
			_connection.PUT(ctx, uri, payload, func(msg []byte, err error) {})
		} else {
			log.WithFields(log.Fields{
				"Error": err.Error(),
			}).Error("Tradfri - SetHex")
		}
	})

}

func SetBlind(deviceID int, position int) {

	_devices.GetDeviceInfo(deviceID, func(device *TradfriDevice, err error) {
		if uri, payload, err := device.SetBlind(position); err == nil {
			ctx, done := context.WithTimeout(context.Background(), 2*time.Second)
			defer done()
			_connection.PUT(ctx, uri, payload, func(msg []byte, err error) {})
		} else {
			log.WithFields(log.Fields{
				"Error": err.Error(),
			}).Error("Tradfri - SetBlind")
		}
	})

}
