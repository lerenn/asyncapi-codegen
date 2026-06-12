package middlewares_test

import (
	"context"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/middlewares"
	"github.com/stretchr/testify/require"
)

type validationPayload struct {
	Name string `json:"name" validate:"required"`
}

func TestValidationAllowsValidPayload(t *testing.T) {
	mw := middlewares.Validation(validator.New(), func() any { return &validationPayload{} })

	called := false
	next := func(_ context.Context) error {
		called = true
		return nil
	}

	msg := &extensions.BrokerMessage{Payload: []byte(`{"name":"alice"}`)}
	err := mw(context.Background(), msg, next)

	require.NoError(t, err)
	require.True(t, called, "next middleware must be called for a valid payload")
}

func TestValidationRejectsInvalidPayload(t *testing.T) {
	mw := middlewares.Validation(validator.New(), func() any { return &validationPayload{} })

	called := false
	next := func(_ context.Context) error {
		called = true
		return nil
	}

	// Missing the required "name" field.
	msg := &extensions.BrokerMessage{Payload: []byte(`{}`)}
	err := mw(context.Background(), msg, next)

	require.Error(t, err)
	require.False(t, called, "next middleware must not be called for an invalid payload")
}
