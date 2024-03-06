package asyncapiv3

import (
	"testing"

	"github.com/lerenn/asyncapi-codegen/pkg/utils"
	"github.com/stretchr/testify/suite"
)

func TestMessageSuite(t *testing.T) {
	suite.Run(t, new(MessageSuite))
}

type MessageSuite struct {
	suite.Suite
}

func (suite *MessageSuite) TestIsCorrelationIDRequiredWithTopLevelHeader() {
	// Set message
	msg := Message{
		Headers: &Schema{
			Required: []string{"correlationId"},
		},
		CorrelationID: &CorrelationID{
			Location: "$message.header#/correlationId",
		},
	}

	// Check if true
	suite.Require().True(msg.isCorrelationIDRequired())
}

func (suite *MessageSuite) TestIsCorrelationIDRequiredWithTopLevelPayload() {
	// Set message
	msg := Message{
		Payload: &Schema{
			Required: []string{"correlationId"},
		},
		CorrelationID: &CorrelationID{
			Location: "$message.payload#/correlationId",
		},
	}

	// Check if true
	suite.Require().True(msg.isCorrelationIDRequired())
}

func (suite *MessageSuite) TestIsCorrelationIDRequiredWithDeepLevelHeader() {
	// Set message
	msg := Message{
		Headers: &Schema{
			Properties: map[string]*Schema{
				"toto": utils.ToPointer(Schema{
					Required: []string{"correlationId"},
				}),
			},
		},
		CorrelationID: &CorrelationID{
			Location: "$message.header#/toto/correlationId",
		},
	}

	// Check if true
	suite.Require().True(msg.isCorrelationIDRequired())
}

func (suite *MessageSuite) TestIsCorrelationIDRequiredWithDeepLevelPayload() {
	// Set message
	msg := Message{
		Payload: &Schema{
			Properties: map[string]*Schema{
				"toto": utils.ToPointer(Schema{
					Required: []string{"correlationId"},
				}),
			},
		},
		CorrelationID: &CorrelationID{
			Location: "$message.payload#/toto/correlationId",
		},
	}

	// Check if true
	suite.Require().True(msg.isCorrelationIDRequired())
}

func (suite *MessageSuite) TestIsCorrelationIDRequiredWithReferencedHeaders() {
	// Set message
	msg := Message{
		Headers: &Schema{
			ReferenceTo: utils.ToPointer(Schema{
				Required: []string{"correlationId"},
			}),
		},
		CorrelationID: &CorrelationID{
			Location: "$message.header#/correlationId",
		},
	}

	// Check if true
	suite.Require().True(msg.isCorrelationIDRequired())
}

func (suite *MessageSuite) TestIsCorrelationIDRequiredWithReferencedPayload() {
	// Set message
	msg := Message{
		Payload: &Schema{
			ReferenceTo: utils.ToPointer(Schema{
				Required: []string{"correlationId"},
			}),
		},
		CorrelationID: &CorrelationID{
			Location: "$message.payload#/correlationId",
		},
	}

	// Check if true
	suite.Require().True(msg.isCorrelationIDRequired())
}

func (suite *MessageSuite) TestIsCorrelationIDRequiredWithInexistantHeaders() {
	// Set message
	msg := Message{
		CorrelationID: &CorrelationID{
			Location: "$message.header#/correlationId",
		},
	}

	// Check if true
	suite.Require().False(msg.isCorrelationIDRequired())
}

func (suite *MessageSuite) TestIsCorrelationIDRequiredWithInexistantPayload() {
	// Set message
	msg := Message{
		CorrelationID: &CorrelationID{
			Location: "$message.payload#/correlationId",
		},
	}

	// Check if true
	suite.Require().False(msg.isCorrelationIDRequired())
}

func (suite *MessageSuite) TestIsCorrelationIDRequiredWithEmptyMessage() {
	// Set message
	msg := Message{}

	// Check if true
	suite.Require().False(msg.isCorrelationIDRequired())
}
