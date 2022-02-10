package tradfri

import (
	"fmt"
	"strings"

	"github.com/buger/jsonparser"
)

type TradfriLight struct {
	Id               int64
	Name             string
	StateDescription string
	State            bool
	Dimmer           int64
	Manufacturer     string
	Model            string
	ColorSpace       string
	ColorValues      ColorDefinition
	Hex              string
}

func (d *TradfriLight) SetState(state int) (uri string, payload string, err error) {
	uri = fmt.Sprintf("%s/%d", uriDevices, d.Id)
	payload = fmt.Sprintf("{ \"%s\": [{ \"%s\": %d }] }", attrLightControl, attrLightState, state)
	return uri, payload, nil
}

func ParseLightInfo(aDevice []byte) (TradfriLight, error) {
	var aLight TradfriLight

	if value, err := jsonparser.GetString(aDevice, attrName); err == nil {
		aLight.Name = value
	}

	if value, err := jsonparser.GetInt(aDevice, attrId); err == nil {
		aLight.Id = value
	}

	_, err := jsonparser.ArrayEach(aDevice, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if res, err := jsonparser.GetInt(value, attrLightState); err == nil {
			aLight.StateDescription, aLight.State = func() (string, bool) {
				if res == 1 {
					return "On", true
				} else {
					return "Off", false
				}
			}()
		}

		if res, err := jsonparser.GetInt(value, attrLightDimmer); err == nil {
			aLight.Dimmer = res
		}

		if res, err := jsonparser.GetString(value, attrLightHex); err == nil {
			aLight.Hex = res
		}

		if res, err := jsonparser.GetInt(value, attrColorX); err == nil {
			aLight.ColorValues.X = res
		} else {
			aLight.ColorValues.X = -1
		}

		if res, err := jsonparser.GetInt(value, attrColorY); err == nil {
			aLight.ColorValues.Y = res
		} else {
			aLight.ColorValues.Y = -1
		}

	}, attrLightControl)
	if err != nil {
		return aLight, err
	}

	if value, err := jsonparser.GetString(aDevice, attrDeviceInfo, attrDeviceInfo_Model); err == nil {
		aLight.Model = value
		if strings.Contains(value, " CWS ") {
			aLight.ColorSpace = "CWS"
		} else if strings.Contains(value, " WS ") {
			aLight.ColorSpace = "WS"
		} else if strings.Contains(value, " WW ") {
			aLight.ColorSpace = "WW"
		} else {
			aLight.ColorSpace = ""
		}
	}

	if value, err := jsonparser.GetString(aDevice, attrDeviceInfo, attrDeviceInfo_Model); err == nil {
		aLight.Model = value
	}

	if value, err := jsonparser.GetString(aDevice, attrDeviceInfo, attrDeviceInfo_Manufacturer); err == nil {
		aLight.Manufacturer = value
	}

	return aLight, err
}
