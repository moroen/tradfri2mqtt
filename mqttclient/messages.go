package mqttclient

type MQTTmessage struct {
	Value int `json:"value"`
}

type MQTTboolMessage struct {
	Value bool `json:"value"`
}

type MQTTLightMessage struct {
	Color struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	} `json:"color"`
	State      string `json:"state"`
	Brightness int    `json:"brightness"`
	ColorTemp  int    `json:"color_temp"`
	ColorMode  string `json:"color_mode"`
}

func SendState(msg []byte) error {
	// fmt.Printf("%s\n", msg)
	/*

		var topic string
		var valueJson []byte

		// SendConfigObject(msg)

		if info, err := coap.ParseLightInfo(msg); err == nil {
			topic = fmt.Sprintf("tradfri/%d/dimmer", info.Id)

			var state string
			if info.State {
				state = "ON"
			} else {
				state = "OFF"
			}

			var color_temp int
			var color_mode string

			switch info.ColorSpace {
			case "WW":
				color_mode = "brightness"
			case "CWS":
				color_mode = "xy"
			case "WS":
				colorMap := coap.CWmap()
				switch info.Hex {
				case colorMap[10]["Hex"]:
					color_temp = 50
				case colorMap[20]["Hex"]:
					color_temp = 320
				case colorMap[30]["Hex"]:
					color_temp = 480
				}
				color_mode = "color_temp"
			}

			message := MQTTLightMessage{State: state, Brightness: int(info.Dimmer), ColorMode: color_mode, ColorTemp: color_temp}
			message.Color.X, message.Color.Y = info.ColorValues.ToFloat()

			if messageJson, err := json.Marshal(message); err == nil {
				log.WithFields(log.Fields{
					"Topic":   topic,
					"Message": string(messageJson),
				}).Debug("Show - Send dimmer message")
				return SendTopic(topic, messageJson, false)
			} else {
				log.WithFields(log.Fields{
					"Error": err.Error(),
				}).Error("Show - Send dimmer message")
			}

			if err != nil {
				log.Fatal(err.Error())
			}
			return err

		} else if info, err := coap.ParsePlugInfo(msg); err == nil {
			topic = fmt.Sprintf("tradfri/%d/switch", info.Id)
			if valueJson, err = json.Marshal(MQTTboolMessage{Value: info.State}); err == nil {
				log.WithFields(log.Fields{
					"Topic":   topic,
					"Message": string(valueJson),
				}).Debug("Show - Send plug message")
				return SendTopic(topic, valueJson, false)
			} else {
				log.Error(err.Error())
				return err
			}

		} else if info, err := coap.ParseBlindInfo(msg); err == nil {
			topic = fmt.Sprintf("tradfri/%d/blind", info.Id)
			if valueJson, err = json.Marshal(MQTTmessage{Value: int(info.Level)}); err == nil {
				log.WithFields(log.Fields{
					"Topic":   topic,
					"Message": string(valueJson),
				}).Debug("Show - Send cover message")
				return SendTopic(topic, valueJson, false)
			} else {
				log.Error(err.Error())
				return err
			}
		} else {

				log.WithFields(log.Fields{
					"Coap message": string(msg),
					"Error":        err.Error(),
				}).Info()

			return err
		}
	*/
	return nil
}
