package tradfri

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

func (d *DeviceList) Init() {
	d.devices = map[int]TradfriDevice{}
}

type TradfriDevice struct {
	DeviceInfo struct {
		Manufacturer    string      `json:"0"`
		Model           string      `json:"1"`
		Num2            string      `json:"2"`
		FirmvareVersion string      `json:"3"`
		Num6            int         `json:"6"`
		Num7            int         `json:"7"`
		Num8            int         `json:"8"`
		BatteryStatus   interface{} `json:"9"`
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
var ErrorNoDimmerControl = errors.New("tradfri Error: device has no dimmer control")
var ErrorNoBlindControl = errors.New("tradfri Error: device has no blind control")
var ErrorNoHexControl = errors.New("tradfri Error: device has no hex control")

type WSStateObject struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Model      string `json:"model"`
	Type       int    `json:"type"`
	ColorSpace string `json:"colorspace"`
	State      bool   `json:"state"`
	Dimmer     int    `json:"dimmer"`
	Position   int    `json:"position"`
	Battery    int    `json:"battery"`
}

func (d *TradfriDevice) WSStateObject() WSStateObject {
	// WS
	ws := WSStateObject{
		Id:         d.Id,
		Name:       d.Name,
		Model:      d.DeviceInfo.Model,
		Type:       d.DeviceType,
		State:      d.GetState(),
		ColorSpace: d.ColorSpace(),
		Dimmer:     d.DimmerLevel(),
		Position:   d.BlindPosition(),
		Battery:    d.BatteryState(),
	}

	return ws
}

func (d *TradfriDevice) DimmerLevel() int {
	if d.LightControl != nil {
		return d.LightControl[0].Dimmer
	} else {
		return -1
	}
}

func (d *TradfriDevice) BlindPosition() int {
	if d.BlindControl != nil {
		return int(d.BlindControl[0].Position)
	} else {
		return -1
	}
}

func (d *TradfriDevice) BatteryState() int {

	if bat, ok := d.DeviceInfo.BatteryStatus.(float64); ok {
		return int(bat)
	} else {
		return -1
	}

}

func (d *TradfriDevice) ColorSpace() string {

	if strings.Contains(d.DeviceInfo.Model, " CWS ") {
		return "CWS"
	} else if strings.Contains(d.DeviceInfo.Model, " WS ") {
		return "WS"
	} else if strings.Contains(d.DeviceInfo.Model, " WW ") {
		return "WW"
	} else {
		return "None"
	}
}

func (d *TradfriDevice) GetType() string {
	if d.LightControl != nil {
		return "dimmer"
	} else if d.PlugControl != nil {
		return "plug"
	} else if d.BlindControl != nil {
		return "blind"
	} else {
		return "NA"
	}
}

func (d *TradfriDevice) GetState() bool {
	if d.LightControl != nil {
		return d.LightControl[0].State != 0
	} else if d.PlugControl != nil {
		return d.PlugControl[0].State != 0
	} else {
		return false
	}
}

// Performs actions
func (d *TradfriDevice) SetState(state bool, handler func([]byte, error)) {
	var value int

	if state {
		value = 1
	} else {
		value = 0
	}

	ctx, done := context.WithTimeout(context.Background(), 10*time.Second)
	defer done()

	if d.LightControl != nil {
		_connection.PUT(ctx, fmt.Sprintf("%s/%d", uriDevices, d.Id), fmt.Sprintf("{ \"%s\": [{ \"%s\": %d }] }", attrLightControl, attrLightState, value), handler)
	} else if d.PlugControl != nil {
		_connection.PUT(ctx, fmt.Sprintf("%s/%d", uriDevices, d.Id), fmt.Sprintf("{ \"%s\": [{ \"%s\": %d }] }", attrPlugControl, attrLightState, value), handler)
	} else {
		handler(nil, ErrorNoStateControl)
	}
}

func (d *TradfriDevice) SetLevel(level int, handler func([]byte, error)) {
	ctx, done := context.WithTimeout(context.Background(), 10*time.Second)
	defer done()

	if d.LightControl != nil {
		_connection.PUT(ctx, fmt.Sprintf("%s/%d", uriDevices, d.Id), fmt.Sprintf("{ \"%s\": [{ \"%s\": %d, \"%s\": %d }] }", attrLightControl, attrLightDimmer, level, attrTransitionTime, 10), handler)
	} else {
		handler(nil, ErrorNoDimmerControl)
	}
}

func (d *TradfriDevice) SetXY(x int, y int, handler func([]byte, error)) {
	ctx, done := context.WithTimeout(context.Background(), 10*time.Second)
	defer done()

	if d.LightControl != nil {
		_connection.PUT(ctx, fmt.Sprintf("%s/%d", uriDevices, d.Id), fmt.Sprintf("{ \"%s\": [{\"%s\": %d, \"%s\": %d}] }", attrLightControl, attrColorX, x, attrColorY, y), handler)
	} else {
		handler(nil, ErrorNoDimmerControl)
	}
}

func (d *TradfriDevice) SetHex(hex string, handler func([]byte, error)) {
	ctx, done := context.WithTimeout(context.Background(), 10*time.Second)
	defer done()

	if d.LightControl != nil {
		_connection.PUT(ctx, fmt.Sprintf("%s/%d", uriDevices, d.Id), fmt.Sprintf("{ \"%s\": [{ \"%s\": \"%s\" }] }", attrLightControl, attrLightHex, hex), handler)
	} else {
		handler(nil, ErrorNoDimmerControl)
	}
}

func (d *TradfriDevice) SetBlind(position int, handler func([]byte, error)) {
	ctx, done := context.WithTimeout(context.Background(), 10*time.Second)
	defer done()

	if d.BlindControl != nil {
		_connection.PUT(ctx, fmt.Sprintf("%s/%d", uriDevices, d.Id), fmt.Sprintf("{ \"%s\": [{ \"%s\": %d}] }", attrBlindControl, attrBlindPosition, position), handler)
	} else {
		handler(nil, ErrorNoBlindControl)
	}
}
