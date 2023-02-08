package asyncapi

import (
	"strings"
)

type Specification struct {
	Version    string             `json:"asyncapi"`
	Info       Info               `json:"info"`
	Channels   map[string]Channel `json:"channels"`
	Components Components         `json:"components"`
}

func (s *Specification) Process() {
	s.setMapsValuesName()

	for name, ch := range s.Channels {
		ch.Process(*s)
		s.Channels[name] = ch
	}
	s.Components.Process(*s)
}

func (s *Specification) setMapsValuesName() {
	for name, ch := range s.Channels {
		ch.Name = name
		s.Channels[name] = ch
	}
}

func (s Specification) ReferenceMessage(ref string) Message {
	name := strings.Split(ref, "/")[3]
	return s.Components.Messages[name]
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
