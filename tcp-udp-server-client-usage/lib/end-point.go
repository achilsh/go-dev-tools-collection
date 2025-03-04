package lib

import "io"

// EndpointConf is the interface implemented by all endpoint configurations.
type EndpointConf interface {
	init(*Node) (Endpoint, error)
}

// Endpoint is an endpoint, which can create Channels.
type Endpoint interface {
	// Conf returns the configuration used to initialize the endpoint
	Conf() EndpointConf
	isEndpoint()
}

// endpointChannelProvider is an endpoint that provides multiple channels.
type endpointChannelProvider interface {
	Endpoint
	close()
	oneChannelAtAtime() bool
	provide() (string, io.ReadWriteCloser, error)
}

