package tradfri

import (
	"context"

	coap "github.com/moroen/gocoap/v5"
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
