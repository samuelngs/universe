package server

import (
	"log"
	"sync"

	"github.com/samuelngs/universe/pkg/crypto"
	"github.com/samuelngs/universe/pkg/uuid"
)

// Option func
type Option func(*Options)

// Options for Secure Shell
type Options struct {
	sync.RWMutex
	// Server Id
	ServerID string
	// NoClientAuth is true if clients are allowed to connect without
	// authenticating.
	NoClientAuth bool
	// Enable password authentication
	PasswordAuthentication bool
	// Enable rsa key authentication
	RSAAuthentication bool
	// Secure Shell server listen addr
	ListenAddr string
	// Secure Shell protocol version
	Protocol int
	// HostKey
	HostKey *crypto.PrivateKey
	// Metadata
	Metadata map[string]string
}

// newOptions creates new option
func newOptions(opts ...Option) *Options {
	o := &Options{
		NoClientAuth:           false,
		PasswordAuthentication: false,
		RSAAuthentication:      false,
		ListenAddr:             ":0",
		Protocol:               2,
		HostKey:                nil,
		Metadata:               make(map[string]string, 0),
		ServerID:               uuid.MustV4(),
	}
	for _, opt := range opts {
		opt(o)
	}
	if o.HostKey == nil {
		k, err := crypto.Generate()
		if err != nil {
			log.Fatal(err)
		}
		o.SetHostKey(k)
	}
	return o
}

// ID option
func ID(s string) Option {
	return func(o *Options) {
		o.SetServerID(s)
	}
}

// ClientAuth option
func ClientAuth(b bool) Option {
	return func(o *Options) {
		o.SetClientAuth(!b)
	}
}

// PasswordAuthentication option
func PasswordAuthentication(b bool) Option {
	return func(o *Options) {
		o.SetPasswordAuthentication(b)
	}
}

// RSAAuthentication option
func RSAAuthentication(b bool) Option {
	return func(o *Options) {
		o.SetRSAAuthentication(b)
	}
}

// ListenAddr option
func ListenAddr(s string) Option {
	return func(o *Options) {
		o.SetListenAddr(s)
	}
}

// Protocol option
func Protocol(v int) Option {
	return func(o *Options) {
		o.SetProtocol(v)
	}
}

// HostKey option
func HostKey(k *crypto.PrivateKey) Option {
	return func(o *Options) {
		o.SetHostKey(k)
	}
}

// Metadata option
func Metadata(m map[string]string) Option {
	return func(o *Options) {
		o.Lock()
		defer o.Unlock()
		for k, v := range m {
			o.Metadata[k] = v
		}
	}
}

// SetClientAuth to enable or disable client authentication [true => enable]
func (v *Options) SetClientAuth(enable bool) *Options {
	v.NoClientAuth = !enable
	return v
}

// SetPasswordAuthentication to enable or disable password authentication
func (v *Options) SetPasswordAuthentication(enable bool) *Options {
	v.PasswordAuthentication = enable
	return v
}

// SetRSAAuthentication to enable or disable rsa key authentication
func (v *Options) SetRSAAuthentication(enable bool) *Options {
	v.RSAAuthentication = enable
	return v
}

// SetListenAddr to set listen port
func (v *Options) SetListenAddr(addr string) *Options {
	if len(addr) > 0 {
		v.ListenAddr = addr
	}
	return v
}

// SetProtocol to set secure shell protocol
func (v *Options) SetProtocol(protocol int) *Options {
	if protocol > 0 {
		v.Protocol = protocol
	}
	return v
}

// SetHostKey to set host private key
func (v *Options) SetHostKey(k *crypto.PrivateKey) *Options {
	v.HostKey = k
	return v
}

// SetServerID to set server reference id
func (v *Options) SetServerID(s string) *Options {
	if len(s) > 0 {
		v.ServerID = s
	}
	return v
}

// GetMetadata to return metadata data
func (v *Options) GetMetadata(k string) string {
	v.RLock()
	defer v.RUnlock()
	return v.Metadata[k]
}
