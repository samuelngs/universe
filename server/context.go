package server

import "net"

// AuthenticationType type
type AuthenticationType int8

// Authentication types
const (
	_ AuthenticationType = iota
	AuthenticationPassword
	AuthenticationPublicKey
	AuthenticationKeyboardInteractive
)

func (v AuthenticationType) String() string {
	switch {
	case v == AuthenticationPassword:
		return "password"
	case v == AuthenticationPublicKey:
		return "rsa"
	case v == AuthenticationKeyboardInteractive:
		return "keyboard-interactive"
	default:
		return "unknown"
	}
}

// Handler for middleware
type Handler func(*Context) error

// Context struct
type Context struct {
	typ          AuthenticationType
	laddr, raddr net.Addr
}

// T returns context handle type
func (v *Context) T() AuthenticationType {
	return v.typ
}

// LocalAddr returns local address
func (v *Context) LocalAddr() net.Addr {
	return v.laddr
}

// RemoteAddr returns remote address
func (v *Context) RemoteAddr() net.Addr {
	return v.raddr
}
