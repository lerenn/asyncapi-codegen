package extensions

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestBrokerSuite(t *testing.T) {
	suite.Run(t, new(BrokerSuite))
}

type BrokerSuite struct {
	suite.Suite
}

func (suite *BrokerSuite) TestIsUninitialized() {
	suite.Require().True(BrokerMessage{}.IsUninitialized())

	suite.Require().False(BrokerMessage{
		Headers: make(map[string][]byte),
	}.IsUninitialized())
}
