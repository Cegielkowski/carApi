package middleware

import "carApi/utils/logger"

// Middleware ...
type Middleware struct {
	logger logger.Logger
}

// NewMiddleware will create new an Middleware object
func NewMiddleware(logger logger.Logger) *Middleware {
	return &Middleware{
		logger: logger,
	}
}
