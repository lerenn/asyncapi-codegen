package middleware

import "context"

// Interface is the type of the function that needs to be implemented to use
// middleware functionnality
type Interface func(ctx context.Context, payload any)
