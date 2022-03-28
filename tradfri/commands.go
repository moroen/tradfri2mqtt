package tradfri

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

/*
func State(deviceID int, state int) {
	go func(deviceID int, state int) {
		_devices.GetDevice(deviceID, func(device *TradfriDevice, err error) {
			if uri, payload, err := device.SetState(state); err == nil {
				ctx, done := context.WithTimeout(context.Background(), 10*time.Second)
				defer done()
				_connection.PUT(ctx, uri, payload, func(msg []byte, err error) {})
			} else {
				log.WithFields(log.Fields{
					"Error": err.Error(),
				}).Error("Tradfri - State")
			}
		})
	}(deviceID, state)
}
*/

func SetXY(deviceID int, x int, y int) {
	go func(deviceID int, x int, y int) {
		_devices.GetDevice(deviceID, func(device *TradfriDevice, err error) {
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
	}(deviceID, x, y)
}

func SetHex(deviceID int, hex string) {
	go func(deviceID int, hex string) {
		_devices.GetDevice(deviceID, func(device *TradfriDevice, err error) {
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
	}(deviceID, hex)
}

func SetBlind(deviceID int, position int) {
	go func(deviceID int, position int) {
		_devices.GetDevice(deviceID, func(device *TradfriDevice, err error) {
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
	}(deviceID, position)
}

func RebootGateway() {
	uri := fmt.Sprintf("%s/%s", attrGateway, attrReboot)
	ctx, done := context.WithTimeout(context.Background(), 2*time.Second)
	defer done()

	log.Debug("Rebooting gateway")

	if _connection.ObserveDone != nil {
		_connection.ObserveDone()
		_connection.ObserveWaitGroup.Wait()
	}

	_connection.POST(ctx, uri, "1", func(msg []byte, err error) {
		fmt.Println(string(msg))
	})
	_connection.Disconnect()
}
