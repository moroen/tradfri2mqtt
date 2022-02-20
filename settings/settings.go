package settings

import (
	"errors"
	"fmt"
	"strings"

	"github.com/kirsle/configdir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var ErrConfigIsDirty = errors.New("config is dirty")

func Init(cpath string) {
	if cpath == "" {
		viper.AddConfigPath(configdir.LocalConfig("tradfri2mqtt"))
	} else {
		viper.AddConfigPath(cpath)
	}

	// File
	viper.SetConfigName("tradfri2mqtt")
	viper.SetConfigType("yaml")

	// Env
	viper.SetEnvPrefix("tradfri2mqtt")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetDefault("mqtt.enable", true)
	viper.SetDefault("mqtt.port", 1883)
	viper.SetDefault("mqtt.host", "127.0.0.1")
	viper.SetDefault("mqtt.discoverytopic", "homeassistant")
	viper.SetDefault("mqtt.commandtopic", "tradfri")

	viper.SetDefault("messages.retrylimit", 5)
	viper.SetDefault("messages.retrydelay", 10)

	viper.SetDefault("tradfri.enable", true)
	viper.SetDefault("tradfri.gateway", "127.0.0.1")
	viper.SetDefault("tradfri.identity", "")
	viper.SetDefault("tradfri.passkey", "")
	viper.SetDefault("tradfri.keepalive", 0)
	viper.SetDefault("tradfri.disconnecttimer", 0)

	viper.SetDefault("interface.enable", true)
	viper.SetDefault("interface.root", "./www")
	viper.SetDefault("interface.port", 8321)
	viper.SetDefault("interface.backloglimit", 10)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Debug("Setting - Config file not found. Creating default")

			if err := viper.SafeWriteConfig(); err != nil {
				fmt.Println(err.Error())
			}
		} else {
			fmt.Println("Other error")
			// Config file was found but another error was produced
		}
	} else {
		viper.AutomaticEnv()
	}
}

type Config struct {
	Messages struct {
		RetryLimit int `json:"retryLimit" yaml:"retrylimit" env:"MESSAGE_RETRY_LIMIT" env-default:"5"`
		RetryDelay int `json:"retryDelay" yaml:"retrydelay" env:"MESSAGE_RETRY_DELAY" env-default:"10"`
	} `json:"messages" yaml:"messages"`
	Mqtt struct {
		Enable         bool   `json:"enable" yaml:"enable" env:"MQTT_ENABLE" env-default:"true"`
		Port           int    `json:"port" yaml:"port" env:"MQTT_BROKER_PORT" env-default:"1883"`
		Host           string `json:"host" yaml:"host" env:"MQTT_BROKER_HOST" env-default:"localhost"`
		DiscoveryTopic string `json:"discoverytopic" yaml:"discoverytopic" env:"MQTT_DISCOVERY_TOPIC" env-default:"homeassistant"`
		CommandTopic   string `json:"commandtopic" yaml:"commandtopic" env:"MQTT_COMMAND_TOPIC" env-default:"tradfri"`
	} `json:"mqtt" yaml:"mqtt"`
	Tradfri struct {
		Enable          bool   `json:"enable" yaml:"enable" env:"TRADFRI_ENABLE" env-default:"true"`
		Gateway         string `json:"gateway" yaml:"gateway" mapstructure:"gateway"`
		Identity        string `json:"identity" yaml:"identity"`
		Passkey         string `json:"passkey" yaml:"passkey"`
		KeepAlive       int    `json:"keepalive" yaml:"keepalive" env-default:"0"`
		DisconnectTimer int    `json:"disconnecttimer" yaml:"disconnecttimer" env-default:"0"`
	} `json:"tradfri" yaml:"tradfri"`
	Interface struct {
		Enable       bool   `json:"enable" yaml:"enable" env:"WWW_ENABLE" env-default:"true"`
		ServerRoot   string `json:"gateway" yaml:"gateway" env:"MQTT_INTERFACE_ROOT" env-default:"./www"`
		BacklogLimit int    `json:"backloglimit" yaml:"backloglimit" env:"MQTT_INTERFACE_BACKLOG_LIMIT" env-default: 10`
	} `json:"interface" yaml:"interface"`
}

/*
var _cfg Config
var configDirs configdir.ConfigDir
var _cfgFile string

func GetConfigFile() (string, error) {
	if _cfgFile != "" {
		return _cfgFile, nil
	}
	configDirs = configdir.New("", "tradfri2mqtt")
	configDirs.LocalPath, _ = filepath.Abs("/config")

	if folder := configDirs.QueryFolderContainsFile(("tradfri2mqtt.yml")); folder != nil {
		file := fmt.Sprintf("%s/%s", folder.Path, "tradfri2mqtt.yml")
		return file, nil
	} else {
		return "", errors.New("config file not found")
	}
}

func GetConfig(force_reload bool) Config {
	if _cfg != (Config{}) && !force_reload {
		return _cfg
	}

	if file, err := GetConfigFile(); err == nil {
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
*/
