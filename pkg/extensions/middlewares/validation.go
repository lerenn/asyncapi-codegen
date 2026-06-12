package middlewares

import (
	"context"
	"encoding/json"

	"github.com/go-playground/validator/v10"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
)

// Validation returns a middleware that validates the message payload against
// the `validate:"..."` struct tags that asyncapi-codegen generates on payload
// types (using github.com/go-playground/validator).
//
// newPayload must return a pointer to a fresh payload value to decode the
// message into, for example:
//
//	middlewares.Validation(validator.New(), func() any { return &UserMessagePayload{} })
//
// The middleware unmarshals msg.Payload (JSON) into that value and validates
// it. If the payload cannot be unmarshaled or fails validation, it returns an
// error, which aborts processing of the message — so invalid messages never
// reach user code (on reception) and are not sent (on publication).
//
// A fresh value is created for every message, so the middleware is safe for
// concurrent use.
func Validation(validate *validator.Validate, newPayload func() any) extensions.Middleware {
	return func(ctx context.Context, msg *extensions.BrokerMessage, next extensions.NextMiddleware) error {
		payload := newPayload()

		if err := json.Unmarshal(msg.Payload, payload); err != nil {
			return err
		}

		if err := validate.Struct(payload); err != nil {
			return err
		}

		return next(ctx)
	}
}
