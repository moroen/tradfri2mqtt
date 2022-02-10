package tradfri

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	coap "github.com/moroen/gocoap/v5"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type GatewayConfig struct {
	Gateway   string
	Identity  string
	Passkey   string
	KeepAlive int
}

func ConnectGateway(ctx context.Context, cfg GatewayConfig, onConnect func(), onDisconnect func(), onCanceled func(), onConnectionFailed func()) (coap.CoapDTLSConnection, error) {
	// SetConfig(cfg)
	coapDTLSConnection := coap.CoapDTLSConnection{}

	coapDTLSConnection.Host = cfg.Gateway
	coapDTLSConnection.Port = 5684
	coapDTLSConnection.Ident = cfg.Identity
	coapDTLSConnection.Key = cfg.Passkey

	coapDTLSConnection.OnConnect = onConnect
	coapDTLSConnection.OnDisconnect = onDisconnect
	coapDTLSConnection.OnCanceled = onCanceled
	coapDTLSConnection.OnConnectionFailed = onConnectionFailed

	coapDTLSConnection.Connect()

	return coapDTLSConnection, nil
}

type pskResponse struct {
	FirmvareVersion string `json:"9029"`
	Key             string `json:"9091"`
}

func GetNewPSK(gateway string, key string, handler func(ident string, key string, err error)) {
	new_ident, uuid_err := uuid.NewV4()

	if uuid_err != nil {
		handler("", "", errors.New("tradfri error: unable to create uuid"))
		return
	}

	payload := fmt.Sprintf("{\"%s\":\"%s\"}", attrIdent, new_ident)
	URI := uriIdent

	coapDTLSConnection := coap.CoapDTLSConnection{}

	coapDTLSConnection.Host = gateway
	coapDTLSConnection.Port = 5684
	coapDTLSConnection.Ident = "Client_identity"
	coapDTLSConnection.Key = key

	coapDTLSConnection.OnConnect = func() {
		var response pskResponse

		log.Debug("GeneratePSK - Connected")
		ctx, done := context.WithTimeout(context.Background(), 2*time.Second)
		defer done()

		coapDTLSConnection.POST(ctx, URI, payload, func(msg []byte, err error) {
			if err := json.Unmarshal(msg, &response); err == nil {
				handler(new_ident.String(), response.Key, nil)
				return
			} else {
				handler("", "", errors.New("tradfri error: unable to unmarshal PSK-response"))
				return
			}
		})
	}

	coapDTLSConnection.OnConnectionFailed = func() {
		log.Info(fmt.Sprintf("Tradfri - GeneratePSK: Unable to connected to gateway at [tcp://%s:%s]", viper.GetString("tradfri.gateway"), "5684"))
	}

	coapDTLSConnection.Connect()
}
