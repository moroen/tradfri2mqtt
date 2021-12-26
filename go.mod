module github.com/moroen/tradfri2mqtt

go 1.16

replace github.com/moroen/go-tradfricoap/async => ../go-tradfricoap/async

replace github.com/moroen/gocoap/v5 => ../gocoap/v5

require (
	github.com/buger/jsonparser v1.1.1
	github.com/eclipse/paho.mqtt.golang v1.3.3
	github.com/gin-contrib/static v0.0.1
	github.com/gin-gonic/gin v1.7.2
	github.com/ilyakaznacheev/cleanenv v1.2.5
	github.com/lucasb-eyer/go-colorful v1.2.0
	github.com/moroen/gocoap/v5 v5.0.0-00010101000000-000000000000
	github.com/shibukawa/configdir v0.0.0-20170330084843-e180dbdc8da0
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.8.1
	gopkg.in/yaml.v2 v2.4.0
)
