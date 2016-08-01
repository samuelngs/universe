package server

import "github.com/samuelngs/universe/pkg/crypto"

// Option for Secure Shell
type Option struct {
	// NoClientAuth is true if clients are allowed to connect without
	// authenticating.
	NoClientAuth bool
	// Enable password authentication
	PasswordAuthentication bool
	// Enable rsa key authentication
	RSAAuthentication bool
	// Secure Shell server listen addr
	Addr string
	// Secure Shell protocol version
	Protocol int
	// HostKey
	HostKey *crypto.PrivateKey
}

// SetClientAuth to enable or disable client authentication [true => enable]
func (v *Option) SetClientAuth(enable bool) *Option {
	v.NoClientAuth = !enable
	return v
}

// SetPasswordAuthentication to enable or disable password authentication
func (v *Option) SetPasswordAuthentication(enable bool) *Option {
	v.PasswordAuthentication = enable
	return v
}

// SetRSAAuthentication to enable or disable rsa key authentication
func (v *Option) SetRSAAuthentication(enable bool) *Option {
	v.RSAAuthentication = enable
	return v
}

// SetAddr to set listen port
func (v *Option) SetAddr(addr string) *Option {
	v.Addr = addr
	return v
}

// SetProtocol to set secure shell protocol
func (v *Option) SetProtocol(protocol int) *Option {
	v.Protocol = protocol
	return v
}
