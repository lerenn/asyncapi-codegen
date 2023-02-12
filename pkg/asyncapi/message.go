package asyncapi

import (
	"strings"

	"github.com/lerenn/asyncapi-codegen/pkg/utils"
)

type Message struct {
	Description   string         `json:"description"`
	Headers       *Any           `json:"headers"`
	Payload       *Any           `json:"payload"`
	CorrelationID *CorrelationID `json:"correlationID"`
	Reference     string         `json:"$ref"`

	// --- Non AsyncAPI fields -------------------------------------------------
	Name string `json:"-"`

	// CorrelationIDLocation will indicate where the correlation id is
	// According to: https://www.asyncapi.com/docs/reference/specification/v2.6.0#correlationIDObject
	CorrelationIDLocation string `json:"-"`
	CorrelationIDRequired bool   `json:"-"`
}

func (msg *Message) Process(spec Specification) {
	msg.CorrelationIDLocation = msg.getCorrelationIDLocation(spec)
	msg.CorrelationIDRequired = msg.isCorrelationIDRequired()
}

func (msg Message) getCorrelationIDLocation(spec Specification) string {
	// Let's check the message before the reference
	if msg.CorrelationID != nil {
		return msg.CorrelationID.Location
	}

	// If there is a reference, check it
	if msg.Reference != "" {
		correlationID := spec.ReferenceMessage(msg.Reference).CorrelationID
		if correlationID != nil {
			return correlationID.Location
		}
	}

	return ""
}

func (msg Message) isCorrelationIDRequired() bool {
	if msg.CorrelationID == nil {
		return false
	}

	path := strings.Split(msg.CorrelationID.Location, "/")

	name := path[len(path)-1]
	if strings.HasPrefix(msg.CorrelationID.Location, "$message.header") && msg.Headers != nil && utils.IsInSlice(msg.Headers.Required, name) {
		return true
	} else if strings.HasPrefix(msg.CorrelationID.Location, "$message.payload") && msg.Payload != nil && utils.IsInSlice(msg.Payload.Required, name) {
		return true
	}

	// TODO: support more sublevels in header and payload
	return false
}
