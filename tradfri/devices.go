package tradfri

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type DeviceList struct {
	mu      sync.Mutex
	devices map[int]TradfriDevice
}

var _devices DeviceList

func (d *DeviceList) UpdateDevice(msg []byte) (*TradfriDevice, error) {
	var device TradfriDevice
	if err := json.Unmarshal(msg, &device); err == nil {
		d.mu.Lock()
		defer d.mu.Unlock()
		d.devices[device.Id] = device
		return &device, nil
	} else {
		return nil, err
	}
}

func (d *DeviceList) GetDevice(id int, forceUpdate bool, handler func(*TradfriDevice, error)) {
	uri := fmt.Sprintf("%s/%d", uriDevices, id)
	ctx, done := context.WithTimeout(context.Background(), 3*time.Second)
	defer done()

	if !forceUpdate {

		if val, ok := d.devices[id]; ok {
			log.WithFields(log.Fields{
				"Cached": true,
			}).Debug("GetDeviceInfo")
			handler(&val, nil)
			return
		}
	}
	_connection.GET(ctx, uri, func(msg []byte, err error) {
		log.WithFields(log.Fields{
			"Cached": false,
		}).Debug("GetDeviceInfo")
		d.mu.Lock()
		defer d.mu.Unlock()

		var device TradfriDevice
		if err := json.Unmarshal(msg, &device); err != nil {
			log.WithFields(log.Fields{
				"Error": err.Error(),
			}).Error("discovery.DeviceList.GetDevices failed")
			handler(nil, err)
		} else {
			d.devices[device.Id] = device
			handler(&device, nil)
		}

	})

}

func GetDevice(id int, forceUpdate bool, handler func(*TradfriDevice, error)) {
	_devices.GetDevice(id, forceUpdate, handler)
}

func GetDevicesList() map[int]TradfriDevice {
	return _devices.devices
}
