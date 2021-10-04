package mqttclient

import (
	"encoding/json"
	"fmt"

	coap "github.com/moroen/go-tradfricoap"
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
	CommandTopic string `json:"command_topic"`
	StateTopic   string `json:"state_topic"`
	// StateValueTemplate      string     `json:"state_value_template"`
	// BrightnessCommandTopic  string     `json:"brightness_command_topic"`
	BrightnessScale int `json:"brightness_scale"`
	// BrightnessStateTopic    string     `json:"brightness_state_topic"`
	// BrightnessValueTemplate string     `json:"brightness_value_template"`
	// OnCommandType       string     `json:"on_command_type"`
	Brightness          string     `json:"brightness"`
	Device              DeviceInfo `json:"device"`
	Name                string     `json:"name"`
	UniqueID            string     `json:"unique_id"`
	ColorMode           string     `json:"color_mode"`
	SupportedColorModes []string   `json:"supported_color_modes"`
}

func SendConfigObject(msg []byte) {

	if light, err := coap.ParseLightInfo(msg); err == nil {
		fmt.Println(light.Name)

		cmdTopic := fmt.Sprintf("tradfri/%d/38/0/dimmer/set", light.Id)
		stTopic := fmt.Sprintf("tradfri/%d/38/0/dimmer", light.Id)
		uniqueID := fmt.Sprintf("tradfri_%d_light", light.Id)
		idents := []string{uniqueID}

		var color_mode string
		var brightness string
		var color_modes []string

		// aConfig := DimmerConfig{StateValueTemplate: "{{ \"OFF\" if value_json.value == 0 else \"ON\" }}", CommandTopic: cmdTopic, StateTopic: stTopic, BrightnessCommandTopic: cmdTopic, BrightnessScale: 99, BrightnessStateTopic: stTopic, BrightnessValueTemplate: "{{ value_json.value }}", OnCommandType: "brightness", Name: light.Name, UniqueID: uniqueID, Device: DeviceInfo{Manufacturer: light.Manufacturer, Identifiers: idents, Model: light.Model, Name: light.Name}}

		switch light.ColorSpace {
		case "WW":
			brightness = "true"
			color_mode = "false"
			color_modes = []string{""}
		case "CWS":
			brightness = "true"
			color_mode = "true"
			color_modes = []string{"xy"}
		case "WS":
			brightness = "true"
			color_mode = "false"
			color_modes = []string{}
		}

		aConfig := DimmerConfig{Brightness: brightness, BrightnessScale: 255, ColorMode: color_mode, SupportedColorModes: color_modes, CommandTopic: cmdTopic, StateTopic: stTopic, Name: light.Name, UniqueID: uniqueID, Device: DeviceInfo{Manufacturer: light.Manufacturer, Identifiers: idents, Model: light.Model, Name: light.Name}}

		// pretty_print(aConfig)
		payload, err := json.Marshal(aConfig)
		if err != nil {
			log.Fatal(err.Error())
		}

		topic := fmt.Sprintf("homeassistant/light/%d/config", light.Id)

		SendTopic(topic, payload)
	} else if plug, err := coap.ParsePlugInfo(msg); err == nil {
		cmdTopic := fmt.Sprintf("tradfri/%d/37/0/switch/set", plug.Id)
		stdTopic := fmt.Sprintf("tradfri/%d/37/0/switch", plug.Id)
		uniqueID := fmt.Sprintf("tradfri_%d_switch", plug.Id)
		idents := []string{uniqueID}
		aConfig := SwitchConfig{PayloadOn: true, PayloadOff: false, ValueTemplate: "{{ value_json.value }}", CommandTopic: cmdTopic, StateTopic: stdTopic, Name: plug.Name, UniqueID: uniqueID, Device: DeviceInfo{Manufacturer: plug.Manufacturer, Identifiers: idents, Model: plug.Model, Name: plug.Name}}

		// pretty_print(aConfig)

		payload, err := json.Marshal(aConfig)
		if err != nil {
			log.Fatal(err.Error())
		}

		SendTopic(fmt.Sprintf("homeassistant/switch/%d/config", plug.Id), payload)

	}

	/*

		} else {
			log.Error(err.Error())
		}

		if plug, err := coap.GetPlug(id); err == nil {
			fmt.Println(plug.Model)

			cmdTopic := fmt.Sprintf("tradfri/%d/37/0/switch/set", id)
			stdTopic := fmt.Sprintf("tradfri/%d/37/0/switch", id)
			uniqueID := fmt.Sprintf("tradfri_%d_switch", id)
			idents := []string{uniqueID}
			aConfig := SwitchConfig{PayloadOn: true, PayloadOff: false, ValueTemplate: "{{ value_json.value }}", CommandTopic: cmdTopic, StateTopic: stdTopic, Name: plug.Name, UniqueID: uniqueID, Device: DeviceInfo{Manufacturer: plug.Manufacturer, Identifiers: idents, Model: plug.Model, Name: plug.Name}}

			pretty_print(aConfig)

			payload, err := json.Marshal(aConfig)
			if err != nil {
				log.Fatal(err.Error())
			}

			SendTopic("homeassistant/switch/65562/switch/config", payload)

		} else {
			log.Error(err.Error())
		}
	*/
}

func pretty_print(val interface{}) {
	if output, err := json.MarshalIndent(val, "", " "); err == nil {
		fmt.Println(string(output))
	}
}
