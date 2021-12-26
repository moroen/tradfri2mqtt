package tradfri

import (
	"fmt"

	"github.com/buger/jsonparser"
)

type TradfriPlug struct {
	Id               int64
	Name             string
	State            bool
	StateDescription string
	Manufacturer     string
	Model            string
}

type TradfriPlugs []TradfriPlug

func (p TradfriPlug) Describe() string {
	return fmt.Sprintf("%d: %s (%s) - %s", p.Id, p.Name, p.Model, p.StateDescription)
}

func ParsePlugInfo(aDevice []byte) (TradfriPlug, error) {
	var p TradfriPlug

	if value, err := jsonparser.GetString(aDevice, attrName); err == nil {
		p.Name = value
	}

	if value, err := jsonparser.GetInt(aDevice, attrId); err == nil {
		p.Id = value
	}

	_, err := jsonparser.ArrayEach(aDevice, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if res, err := jsonparser.GetInt(value, attrPlugState); err == nil {
			p.StateDescription, p.State = func() (string, bool) {
				if res == 1 {
					return "On", true
				} else {
					return "Off", false
				}
			}()
		}
	}, attrPlugControl)
	if err != nil {
		return p, err
	}

	if value, err := jsonparser.GetString(aDevice, attrDeviceInfo, attrDeviceInfo_Model); err == nil {
		p.Model = value
	}

	if value, err := jsonparser.GetString(aDevice, attrDeviceInfo, attrDeviceInfo_Manufacturer); err == nil {
		p.Manufacturer = value
	}

	return p, err
}
