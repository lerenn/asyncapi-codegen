package asyncapi

import (
	"strings"
)

type Specification struct {
	Version    string              `json:"asyncapi"`
	Info       Info                `json:"info"`
	Channels   map[string]*Channel `json:"channels"`
	Components Components          `json:"components"`
}

func (s *Specification) Process() {
	for path, ch := range s.Channels {
		ch.Process(path, *s)
	}
	s.Components.Process(*s)
}

func (s Specification) GetPublishSubscribeCount() (publishCount, subscribeCount uint) {
	for _, c := range s.Channels {
		if c.Publish != nil {
			publishCount++
		} else if c.Subscribe != nil {
			subscribeCount++
		}
	}

	return publishCount, subscribeCount
}

func (s Specification) ReferenceMessage(ref string) *Message {
	msg, _ := s.reference(ref).(*Message)
	return msg
}

func (s Specification) ReferenceAny(ref string) *Any {
	msg, _ := s.reference(ref).(*Any)
	return msg
}

func (s Specification) reference(ref string) interface{} {
	refPath := strings.Split(ref, "/")[1:]

	if refPath[0] == "components" {
		if refPath[1] == "messages" {
			msg := s.Components.Messages[refPath[2]]
			return msg.referenceFrom(refPath[3:])
		} else if refPath[1] == "schemas" {
			schema := s.Components.Schemas[refPath[2]]
			return schema.referenceFrom(refPath[3:])
		}
	}

	return nil
}
