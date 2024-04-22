package test

import (
	"os"
)

// BrokerAddressParams is the parameters for the BrokerAddress function.
type BrokerAddressParams struct {
	Schema string
	Port   string

	DockerizedAddr string
	DockerizedPort string

	LocalAddr string
	LocalPort string
}

// BrokerAddress returns the broker address based on the environment.
// If the environment variable ASYNCAPI_DOCKERIZED is set, it returns
// the dockerized address.
func BrokerAddress(params BrokerAddressParams) string {
	var url string

	// Set schema if not empty
	if params.Schema != "" {
		url = params.Schema + "://"
	}

	// Set address based on environment
	dockerized := (os.Getenv("ASYNCAPI_DOCKERIZED") != "")
	switch {
	case dockerized:
		url += params.DockerizedAddr
	case params.LocalAddr != "":
		url += params.LocalAddr
	default:
		url += "localhost"
	}

	// Set port if not empty
	switch {
	case dockerized && params.DockerizedPort != "":
		url += ":" + params.DockerizedPort
	case !dockerized && params.LocalPort != "":
		url += ":" + params.LocalPort
	case params.Port != "":
		url += ":" + params.Port
	}

	return url
}
