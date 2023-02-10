// Package "generated" provides primitives to interact with the AsyncAPI specification.
//
// Code generated by github.com/lerenn/asyncapi-codegen version (devel) DO NOT EDIT.
package generated

// ClientController is the structure that provides publishing capabilities to the
// developer and and connect the broker with the client
type ClientController struct {
	brokerController BrokerController
	stopSubscribers  map[string]chan interface{}
}

// NewClientController links the client to the broker
func NewClientController(bs BrokerController) (*ClientController, error) {
	if bs == nil {
		return nil, ErrNilBrokerController
	}

	return &ClientController{
		brokerController: bs,
		stopSubscribers:  make(map[string]chan interface{}),
	}, nil
}

// Close will clean up any existing resources on the controller
func (cc *ClientController) Close() {
	// Nothing to do
}

// PublishHello will publish messages to 'hello' channel
func (cc *ClientController) PublishHello(msg HelloMessage) error {
	// TODO: check that 'cc' is not nil

	// Convert to UniversalMessage
	um, err := msg.toUniversalMessage()
	if err != nil {
		return err
	}

	// Publish on event broker
	return cc.brokerController.Publish("hello", um)
}

// Listen will let the controller handle subscriptions and will be interrupted
// only when an struct is sent on the interrupt channel
func (cc *ClientController) Listen(irq chan interface{}) {
	<-irq
}
