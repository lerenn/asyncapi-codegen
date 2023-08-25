package middleware

import "context"

// Interface is the type of the function that needs to be implemented to use
// middleware functionnality
//
// As payload is a generic type, it is up to the user to cast it to the correct type
// but the user can also alter the payload
type Interface func(ctx context.Context, payload any) context.Context
