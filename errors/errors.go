package errors

import "errors"

// ErrorTimeout error
var ErrorBrokerNotConnected = errors.New("Broker Error: Connection refused")

// ErrorTimeout error
var ErrorConfigStale = errors.New("Settings Error: Stale configuration")
