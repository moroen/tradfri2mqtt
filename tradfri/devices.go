package tradfri

import (
	"context"
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type DeviceList struct {
	mu      sync.Mutex
	devices map[int64][]byte
}

func (d *DeviceList) Init() {
	d.devices = map[int64][]byte{}
}

var _devices DeviceList

func (d *DeviceList) GetDeviceInfo(id int64, handler func([]byte, error)) {
	uri := fmt.Sprintf("%s/%d", uriDevices, id)
	ctx, done := context.WithTimeout(context.Background(), 3*time.Second)
	defer done()

	if val, ok := d.devices[id]; ok {
		log.WithFields(log.Fields{
			"Cached": true,
		}).Debug("GetDeviceInfo")
		handler(val, nil)
		return
	}

	_connection.GET(ctx, uri, func(msg []byte, err error) {
		log.WithFields(log.Fields{
			"Cached": false,
		}).Debug("GetDeviceInfo")
		d.mu.Lock()
		d.devices[id] = msg
		d.mu.Unlock()
		handler(msg, nil)
	})
}
