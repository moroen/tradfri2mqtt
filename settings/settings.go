package settings

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/ilyakaznacheev/cleanenv"
	coap "github.com/moroen/go-tradfricoap"
	"github.com/shibukawa/configdir"
	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Messages struct {
		RetryLimit int `json:"retryLimit" yaml:"retrylimit" env:"MESSAGE_RETRY_LIMIT" env-default:"5"`
		RetryDelay int `json:"retryDelay" yaml:"retrydelay" env:"MESSAGE_RETRY_DELAY" env-default:"10"`
	} `json:"messages" yaml:"messages"`
	Mqtt struct {
		Port string `json:"port" yaml:"port" env:"MQTT_BROKER_PORT" env-default:"1883"`
		Host string `json:"host" yaml:"host" env:"MQTT_BROKER_HOST" env-default:"localhost"`
	} `json:"mqtt" yaml:"mqtt"`
	Tradfri struct {
		Gateway   string `json:"gateway" yaml:"gateway"`
		Identity  string `json:"ident" yaml:"ident"`
		Passkey   string `json:"key" yaml:"key"`
		KeepAlive int    `json:"keepAlive" yaml:"keepalive" env-default:"0"`
	} `json:"tradfri" yaml:"tradfri"`
	Interface struct {
		ServerRoot string `json:"gateway" yaml:"gateway" env:"MQTT_INTERFACE_ROOT" env-default:"./www"`
	} `json:"interface" yaml:"interface"`
}

var _cfg Config
var configDirs configdir.ConfigDir

func GetConfig(force_reload bool) Config {
	if _cfg != (Config{}) && !force_reload {
		return _cfg
	}

	configDirs = configdir.New("", "tradfri2mqtt")
	configDirs.LocalPath, _ = filepath.Abs("/config")

	if folder := configDirs.QueryFolderContainsFile(("tradfri2mqtt.yml")); folder != nil {
		file := fmt.Sprintf("%s/%s", folder.Path, "tradfri2mqtt.yml")
		log.WithFields(log.Fields{
			"File": file,
		}).Debug("Loading config")
		cleanenv.ReadConfig(file, &_cfg)
	} else {
		log.Debug("No config found")
		cleanenv.ReadEnv(&_cfg)
		WriteConfig(&_cfg)
	}

	return _cfg
}

func (c Config) Save() {
	WriteConfig(&c)
}

func WriteConfig(cfg *Config) (err error) {
	s, err := yaml.Marshal(cfg)
	if err != nil {
		panic(err.Error())
	}

	var folders []*configdir.Config

	if _, err := os.Stat("/config"); os.IsNotExist(err) {
		folders = configDirs.QueryFolders(configdir.Global)
	} else {
		folders = configDirs.QueryFolders(configdir.Local)
	}

	if err = folders[0].WriteFile("tradfri2mqtt.yml", s); err != nil {
		log.WithFields(log.Fields{
			"Folder": folders[0].Path,
		}).Error("Error saving config")
	} else {
		log.WithFields(log.Fields{
			"File": fmt.Sprintf("%s/tradfri2mqtt.yml", folders[0].Path),
		}).Debug("Saving config")
	}

	return err
}

func GetCoapConfig(force_reload bool) coap.GatewayConfig {
	cfg := GetConfig(force_reload)

	coapSettings := coap.GatewayConfig{
		Gateway:   cfg.Tradfri.Gateway,
		Identity:  cfg.Tradfri.Identity,
		Passkey:   cfg.Tradfri.Passkey,
		KeepAlive: cfg.Tradfri.KeepAlive,
	}
	return coapSettings
}
