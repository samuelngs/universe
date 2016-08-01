package server

import "github.com/samuelngs/universe/errors"

const namespace string = "server"

// Error messages
var (
	ErrUnauthentized = errors.Unauthorized(namespace, "connection unauthorized")
)
