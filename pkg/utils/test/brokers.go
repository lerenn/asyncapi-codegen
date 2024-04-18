package test

import (
	"os"
)

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
	if os.Getenv("ASYNCAPI_DOCKERIZED") != "" {
		url += params.DockerizedAddr

		// Set port if not empty
		if params.DockerizedPort != "" {
			url += ":" + params.DockerizedPort
		} else if params.Port != "" {
			url += ":" + params.Port
		}
	} else {
		if params.LocalAddr != "" {
			url += params.LocalAddr
		} else {
			url += "localhost"
		}

		// Set port if not empty
		if params.LocalPort != "" {
			url += ":" + params.LocalPort
		} else if params.Port != "" {
			url += ":" + params.Port
		}
	}

	return url
}
