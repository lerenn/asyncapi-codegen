//go:generate go run ../../../../cmd/asyncapi-codegen -p issue49 -i ./asyncapi.yaml -o ./asyncapi.gen.go

package issue49

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

type Suite struct {
	suite.Suite
}

func (suite *Suite) TestCorrectPublicationsSubscriptionsAreGenerated() {
	// Check that the User subscriber is generated with correct subscriptions
	_ = UserSubscriber.Issue49Chat
	_ = UserSubscriber.Issue49Status

	// Check that the User publisher is generated with correct publications
	userController := UserController{}
	_ = userController.PublishIssue49Chat

	// Check that the User subscriber is generated with correct subscriptions
	_ = AppSubscriber.Issue49Chat

	// Check that the App publisher is generated with correct publications
	appController := AppController{}
	_ = appController.PublishIssue49Status
	_ = appController.PublishIssue49Chat
}
