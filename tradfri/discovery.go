package tradfri

import (
	"encoding/json"
	"fmt"

	"github.com/buger/jsonparser"
	log "github.com/sirupsen/logrus"
)

type SwitchConfig struct {
	PayloadOff    bool       `json:"payload_off"`
	PayloadOn     bool       `json:"payload_on"`
	ValueTemplate string     `json:"value_template"`
	CommandTopic  string     `json:"command_topic"`
	StateTopic    string     `json:"state_topic"`
	Device        DeviceInfo `json:"device"`
	Name          string     `json:"name"`
	UniqueID      string     `json:"unique_id"`
}

type DeviceInfo struct {
	Identifiers  []string `json:"identifiers"`
	Manufacturer string   `json:"manufacturer"`
	Model        string   `json:"model"`
	Name         string   `json:"name"`
	SwVersion    string   `json:"sw_version"`
}

type DimmerConfig struct {
	Schema       string `json:"schema"`
	CommandTopic string `json:"command_topic"`
	StateTopic   string `json:"state_topic"`
	// StateValueTemplate      string     `json:"state_value_template"`
	// BrightnessCommandTopic  string     `json:"brightness_command_topic"`
	BrightnessScale int `json:"brightness_scale"`
	// BrightnessStateTopic    string     `json:"brightness_state_topic"`
	// BrightnessValueTemplate string     `json:"brightness_value_template"`
	// OnCommandType       string     `json:"on_command_type"`
	Brightness          bool       `json:"brightness"`
	Device              DeviceInfo `json:"device"`
	Name                string     `json:"name"`
	UniqueID            string     `json:"unique_id"`
	ColorMode           bool       `json:"color_mode"`
	SupportedColorModes []string   `json:"supported_color_modes"`
}

type BlindConfig struct {
	CommandTopic        string     `json:"command_topic"`
	PositionTopic       string     `json:"position_topic"`
	PositionTemplate    string     `json:"position_template"`
	SetPositionTopic    string     `json:"set_position_topic"`
	SetPositionTemplate string     `json:"set_position_template"`
	PayloadOpen         string     `json:"payload_open"`
	PayloadClose        string     `json:"payload_close"`
	PayloadStop         string     `json:"payload_stop"`
	Device              DeviceInfo `json:"device"`
	Name                string     `json:"name"`
	UniqueID            string     `json:"unique_id"`
}

var _discovered bool

func Discover(force bool) {
	if force || !_discovered {
		_connection.GET(_ctx, uriDevices, func(msg []byte, err error) {
			fmt.Println(string(msg))
			_, err = jsonparser.ArrayEach(msg, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
				if res, err := jsonparser.GetInt(value); err == nil {
					uri := fmt.Sprintf("%s/%d", uriDevices, res)
					_connection.GET(_ctx, uri, func(msg []byte, err error) {
						SendConfigObject(msg)
					})
				} else {
					fmt.Println(err.Error())
				}
			})
		})
		_discovered = true
	}
}

func SendConfigObject(msg []byte) {
	if light, err := ParseLightInfo(msg); err == nil {
		// fmt.Println(light.Name)

		cmdTopic := fmt.Sprintf("%s/%d/dimmer/set", _mqtt_command_topic, light.Id)
		stTopic := fmt.Sprintf("%s/%d/dimmer", _mqtt_command_topic, light.Id)
		uniqueID := fmt.Sprintf("tradfri_%d_light", light.Id)
		idents := []string{uniqueID}

		var color_mode bool
		var brightness bool
		var color_modes []string

		// aConfig := DimmerConfig{StateValueTemplate: "{{ \"OFF\" if value_json.value == 0 else \"ON\" }}", CommandTopic: cmdTopic, StateTopic: stTopic, BrightnessCommandTopic: cmdTopic, BrightnessScale: 99, BrightnessStateTopic: stTopic, BrightnessValueTemplate: "{{ value_json.value }}", OnCommandType: "brightness", Name: light.Name, UniqueID: uniqueID, Device: DeviceInfo{Manufacturer: light.Manufacturer, Identifiers: idents, Model: light.Model, Name: light.Name}}

		switch light.ColorSpace {
		case "WW":
			brightness = true
			color_mode = false
			color_modes = []string{"brightness"}
		case "CWS":
			brightness = true
			color_mode = true
			color_modes = []string{"xy"}
		case "WS":
			brightness = true
			color_mode = false
			color_modes = []string{"color_temp"}
		}

		aConfig := DimmerConfig{
			Schema:              "json",
			Brightness:          brightness,
			BrightnessScale:     255,
			ColorMode:           color_mode,
			SupportedColorModes: color_modes,
			CommandTopic:        cmdTopic,
			StateTopic:          stTopic,
			Name:                light.Name,
			UniqueID:            uniqueID,
			Device: DeviceInfo{
				Manufacturer: light.Manufacturer,
				Identifiers:  idents,
				Model:        light.Model,
				Name:         light.Name,
			},
		}

		// pretty_print(aConfig)
		payload, err := json.Marshal(aConfig)
		if err != nil {
			log.Fatal(err.Error())
		}

		topic := fmt.Sprintf("%s/light/%d/config", _mqtt_discovery_topic, light.Id)

		log.WithFields(log.Fields{
			"topic":  topic,
			"config": string(payload),
		}).Debug("Disovery - Light")
		MQTTSendTopic(topic, payload, true)

	} else if plug, err := ParsePlugInfo(msg); err == nil {

		cmdTopic := fmt.Sprintf("%s/%d/switch/set", _mqtt_command_topic, plug.Id)
		stdTopic := fmt.Sprintf("%s/%d/switch", _mqtt_command_topic, plug.Id)
		uniqueID := fmt.Sprintf("tradfri_%d_switch", plug.Id)
		idents := []string{uniqueID}
		aConfig := SwitchConfig{
			PayloadOn:     true,
			PayloadOff:    false,
			ValueTemplate: "{{ value_json.value }}",
			CommandTopic:  cmdTopic,
			StateTopic:    stdTopic,
			Name:          plug.Name,
			UniqueID:      uniqueID,
			Device: DeviceInfo{
				Manufacturer: plug.Manufacturer,
				Identifiers:  idents,
				Model:        plug.Model,
				Name:         plug.Name,
			},
		}

		// pretty_print(aConfig)

		payload, err := json.Marshal(aConfig)
		if err != nil {
			log.Fatal(err.Error())
		}

		topic := fmt.Sprintf("%s/switch/%d/config", _mqtt_discovery_topic, plug.Id)

		log.WithFields(log.Fields{
			"topic":  topic,
			"config": string(payload),
		}).Debug("Disovery - Plug")

		MQTTSendTopic(topic, payload, true)

	} else if blind, err := ParseBlindInfo(msg); err == nil {

		cmdTopic := fmt.Sprintf("%s/%d/blind/set", _mqtt_command_topic, blind.Id)
		posTopic := fmt.Sprintf("%s/%d/blind", _mqtt_command_topic, blind.Id)
		setPosTopic := fmt.Sprintf("%s/%d/blind/set", _mqtt_command_topic, blind.Id)
		uniqueID := fmt.Sprintf("tradfri_%d_blind", blind.Id)
		idents := []string{uniqueID}

		aConfig := BlindConfig{
			CommandTopic:        cmdTopic,
			PositionTopic:       posTopic,
			PositionTemplate:    "{{ value_json.position }}",
			SetPositionTopic:    setPosTopic,
			SetPositionTemplate: "{ \"position\": {{ position }} }",
			PayloadOpen:         "{ \"position\": 0 }",
			PayloadClose:        "{ \"position\": 100 }",
			PayloadStop:         "",
			UniqueID:            uniqueID,
			Device:              DeviceInfo{Manufacturer: blind.Manufacturer, Identifiers: idents, Model: blind.Model, Name: blind.Name},
			Name:                blind.Name,
		}

		payload, err := json.Marshal(aConfig)
		if err != nil {
			log.Fatal(err.Error())
		}

		topic := fmt.Sprintf("%s/cover/%d/config", _mqtt_discovery_topic, blind.Id)

		log.WithFields(log.Fields{
			"topic":  topic,
			"config": string(payload),
		}).Debug("Disovery - Blind")

		MQTTSendTopic(topic, payload, true)
	}
}

func pretty_print(val interface{}) {
	if output, err := json.MarshalIndent(val, "", " "); err == nil {
		fmt.Println(string(output))
	}
}
