package tradfri

import (
	"errors"
	"fmt"
	"sort"
)

type ColorMap map[int]map[string]string

func CWmap() ColorMap {
	var whiteBalance = ColorMap{
		0:  {"Name": "Off", "Hex": "000000"},
		10: {"Name": "Cold", "Hex": "f5faf6"},
		20: {"Name": "Normal", "Hex": "f1e0b5"},
		30: {"Name": "Warm", "Hex": "efd275"},
	}

	return whiteBalance
}

func CWSmap() ColorMap {
	var cws = ColorMap{
		0:   {"Name": "Off", "Hex": "000000"},
		10:  {"Name": "Blue", "Hex": "4a418a"},
		20:  {"Name": "Candlelight", "Hex": "ebb63e"},
		30:  {"Name": "Cold sky", "Hex": "dcf0f8"},
		40:  {"Name": "Cool daylight", "Hex": "eaf6fb"},
		50:  {"Name": "Cool white", "Hex": "f5faf6"},
		60:  {"Name": "Dark Peach", "Hex": "da5d41"},
		70:  {"Name": "Light Blue", "Hex": "6c83ba"},
		80:  {"Name": "Light Pink", "Hex": "e8bedd"},
		90:  {"Name": "Light Purple", "Hex": "c984bb"},
		100: {"Name": "Lime", "Hex": "a9d62b"},
		110: {"Name": "Peach", "Hex": "e57345"},
		120: {"Name": "Pink", "Hex": "e491af"},
		130: {"Name": "Saturated Red", "Hex": "dc4b31"},
		140: {"Name": "Saturated Pink", "Hex": "d9337c"},
		150: {"Name": "Saturated Purple", "Hex": "8f2686"},
		160: {"Name": "Sunrise", "Hex": "f2eccf"},
		170: {"Name": "Yellow", "Hex": "d6e44b"},
		180: {"Name": "Warm Amber", "Hex": "e78834"},
		190: {"Name": "Warm glow", "Hex": "efd275"},
		200: {"Name": "Warm white", "Hex": "f1e0b5"},
	}
	return cws
}

type ColorDefinition struct {
	X int64
	Y int64
}

func (c *ColorDefinition) ToFloat() (float64, float64) {
	return float64(c.X) / 65535, float64(c.Y) / 65535
}

func hexForLevel(colorMap ColorMap, level int) string {
	return colorMap[level]["Hex"]
}

func ListColorsInMap(colorMap ColorMap) {
	var keys []int

	for k := range colorMap {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, val := range keys {
		fmt.Printf("%d: %s\n", val, colorMap[val]["Name"])
	}
}

/*
func SetXY(id int64, x int64, y int64) error {
	uri := fmt.Sprintf("%s/%d", uriDevices, 65554)
	payload := fmt.Sprintf("{ \"3311\": [{\"%s\": %d, \"%s\": %d}] }", attrColorX, x, attrColorY, y)

	log.WithFields(log.Fields{
		uri:     uri,
		payload: payload,
	}).Debug()


		if _, err := PutRequest(uri, payload); err != nil {
			return err
		}

	return nil

}
*/

/*
func SetRGB(id int64, rgb string) error {
	fmt.Printf("Device: %d - RGB: %s\n", id, rgb)
	c, err := colorful.Hex(rgb)
	if err != nil {
		log.Fatal(err)
	}

	// h, s, v := c.Hsv()
	// x, y, z := c.Xyz()
	x, y, _ := c.Xyy()

	// fmt.Println("HSV: ", h, s, v)
	// fmt.Println("xyY:", x, y, lum)

	x = x * 65535
	y = y * 65535

	// uri := fmt.Sprintf("%s/%d", uriDevices, 65554)
	payload := fmt.Sprintf("{ \"3311\": [{\"5709\": %d, \"5710\": %d}] }", int(x), int(y))

	//payload := fmt.Sprintf("{ \"%s\": [{ \"%s\": %d }] }", attr_Light_control, attr_light_state, state)
	return nil
}
*/

func GetColorMap(ColorSpace string) (ColorMap, error) {
	switch ColorSpace {
	case "CWS":
		return CWSmap(), nil
	case "WS":
		return CWmap(), nil
	default:
		return nil, errors.New("unknown colorspace")
	}
}

func GetLevelForHex(ColorSpace, hex string) (int, error) {
	if colorMap, err := GetColorMap(ColorSpace); err == nil {
		for key, value := range colorMap {
			if value["Hex"] == hex {
				return key, nil
			}
		}
	}

	return 0, nil
}
