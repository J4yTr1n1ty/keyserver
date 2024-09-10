package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

// CreateStack creates a new stack of middleware functions to avoid ugly notation in the main function
// e.g. stack := middleware.CreateStack(middleware.Logging, middleware.Auth)
func CreateStack(xs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(xs) - 1; i >= 0; i-- {
			x := xs[i]
			next = x(next)
		}

		return next
	}
}
