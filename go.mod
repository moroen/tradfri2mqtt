module github.com/moroen/tradfri2mqtt

go 1.16

replace github.com/moroen/go-tradfricoap => ../go-tradfricoap

replace github.com/moroen/gocoap/v4 => ../gocoap/v4

require (
	github.com/eclipse/paho.mqtt.golang v1.3.3
	github.com/ilyakaznacheev/cleanenv v1.2.5
	github.com/moroen/go-tradfricoap v0.1.1
	github.com/sirupsen/logrus v1.4.2
	github.com/thedevsaddam/retry v0.0.0-20200324223450-9769a859cc6d // indirect
	gopkg.in/yaml.v2 v2.4.0
)
