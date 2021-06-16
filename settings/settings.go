package settings

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/ilyakaznacheev/cleanenv"
	coap "github.com/moroen/go-tradfricoap"
	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Mqtt struct {
		Port string `yaml:"port" env:"MQTT_BROKER_PORT" env-default:"1883"`
		Host string `yaml:"host" env:"MQTT_BROKER_HOST" env-default:"localhost"`
	} `yaml:"mqtt"`
	Tradfri struct {
		Gateway  string `yaml:"gateway"`
		Identity string `yaml:"ident"`
		Passkey  string `yaml:"key"`
	} `yaml:"tradfri"`
}

var _cfg Config

func GetConfig() Config {
	if _cfg != (Config{}) {
		return _cfg
	}

	err := cleanenv.ReadConfig("config/tradfri2mqtt.yml", &_cfg)
	if err != nil {
		err = cleanenv.ReadEnv(&_cfg)
		if err != nil {
			panic(err.Error())
		}
		WriteConfig(&_cfg)
	}
	err = cleanenv.ReadEnv(&_cfg)
	if err != nil {
		panic(err.Error())
	}

	return _cfg
}

func WriteConfig(cfg *Config) (err error) {
	s, err := yaml.Marshal(cfg)
	if err != nil {
		panic(err.Error())
	}
	log.Println("Saving config")

	if _, err := os.Stat("config/"); os.IsNotExist(err) {
		os.Mkdir("config", 0755)
	}

	err = os.WriteFile("config/tradfri2mqtt.yml", s, 0644)
	if err != nil {
		panic(err.Error())
	}
	return err
}

func GetCoapConfig() coap.GatewayConfig {
	cfg := GetConfig()

	coapSettings := coap.GatewayConfig{
		Gateway:  cfg.Tradfri.Gateway,
		Identity: cfg.Tradfri.Identity,
		Passkey:  cfg.Tradfri.Passkey,
	}
	return coapSettings
}
