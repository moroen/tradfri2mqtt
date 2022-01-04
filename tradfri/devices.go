package tradfri

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type DeviceList struct {
	mu      sync.Mutex
	devices map[int]TradfriDevice
}

func (d *DeviceList) Init() {
	d.devices = map[int]TradfriDevice{}
}

type TradfriDevice struct {
	DeviceInfo struct {
		Manufacturer    string `json:"0"`
		Model           string `json:"1"`
		Num2            string `json:"2"`
		FirmvareVersion string `json:"3"`
		Num6            int    `json:"6"`
		Num7            int    `json:"7"`
		Num8            int    `json:"8"`
		BatteryStatus   int    `json:"9"`
	} `json:"3"`
	LightControl []struct {
		ColorHex string `json:"5706"`
		Num5707  int    `json:"5707"`
		Num5708  int    `json:"5708"`
		ColorX   int    `json:"5709"`
		ColorY   int    `json:"5710"`
		Num5849  int    `json:"5849"`
		State    int    `json:"5850"`
		Dimmer   int    `json:"5851"`
		Num9003  int    `json:"9003"`
	} `json:"3311"`
	PlugControl []struct {
		Num5849 int `json:"5849"`
		State   int `json:"5850"`
		Num5851 int `json:"5851"`
		Num9003 int `json:"9003"`
	} `json:"3312"`
	BlindControl []struct {
		Position float64 `json:"5536"`
		Num9003  int     `json:"9003"`
	} `json:"15015"`
	DeviceType int    `json:"5750"`
	Name       string `json:"9001"`
	Num9002    int    `json:"9002"`
	Id         int    `json:"9003"`
	Num9019    int    `json:"9019"`
	Num9020    int    `json:"9020"`
	Num9054    int    `json:"9054"`
}

var ErrorNoStateControl = errors.New("tradfri Error: device has no state control")

var ErrorNoDimmerControl = errors.New("tradfri Error: device has no state control")

func (d *TradfriDevice) SetState(state int) (string, string, error) {
	if d.LightControl != nil {
		return fmt.Sprintf("%s/%d", uriDevices, d.Id), fmt.Sprintf("{ \"%s\": [{ \"%s\": %d }] }", attrLightControl, attrLightState, state), nil
	} else if d.PlugControl != nil {
		fmt.Println("PlugControl")
		return fmt.Sprintf("%s/%d", uriDevices, d.Id), fmt.Sprintf("{ \"%s\": [{ \"%s\": %d }] }", attrPlugControl, attrLightState, state), nil
	} else {
		return "", "", ErrorNoStateControl
	}
}

func (d *TradfriDevice) SetLevel(level int) (string, string, error) {
	if d.LightControl != nil {
		uri := fmt.Sprintf("%s/%d", uriDevices, d.Id)
		payload := fmt.Sprintf("{ \"%s\": [{ \"%s\": %d, \"%s\": %d }] }", attrLightControl, attrLightDimmer, level, attrTransitionTime, 10)
		return uri, payload, nil
	} else {
		return "", "", ErrorNoDimmerControl
	}
}

func (d *TradfriDevice) SetXY(x int, y int) (string, string, error) {
	if d.LightControl != nil {
		uri := fmt.Sprintf("%s/%d", uriDevices, d.Id)
		payload := fmt.Sprintf("{ \"%s\": [{\"%s\": %d, \"%s\": %d}] }", attrLightControl, attrColorX, x, attrColorY, y)
		return uri, payload, nil
	} else {
		return "", "", ErrorNoDimmerControl
	}
}

func (d *TradfriDevice) SetHex(hex string) (string, string, error) {
	if d.LightControl != nil {
		uri := fmt.Sprintf("%s/%d", uriDevices, d.Id)
		payload := fmt.Sprintf("{ \"%s\": [{ \"%s\": \"%s\" }] }", attrLightControl, attrLightHex, hex)
		return uri, payload, nil
	} else {
		return "", "", ErrorNoDimmerControl
	}
}

var _devices DeviceList

func (d *DeviceList) GetDeviceInfo(id int, handler func(*TradfriDevice, error)) {
	uri := fmt.Sprintf("%s/%d", uriDevices, id)
	ctx, done := context.WithTimeout(context.Background(), 3*time.Second)
	defer done()

	if val, ok := d.devices[id]; ok {
		log.WithFields(log.Fields{
			"Cached": true,
		}).Debug("GetDeviceInfo")
		handler(&val, nil)
		return
	}

	_connection.GET(ctx, uri, func(msg []byte, err error) {
		log.WithFields(log.Fields{
			"Cached": false,
		}).Debug("GetDeviceInfo")
		d.mu.Lock()
		var device TradfriDevice
		if err := json.Unmarshal(msg, &device); err != nil {
			fmt.Println(err.Error())
			handler(nil, err)
		} else {
			d.devices[device.Id] = device
			d.mu.Unlock()
			handler(&device, nil)
		}

	})
}
