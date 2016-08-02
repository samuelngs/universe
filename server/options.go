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
	// HostKeys
	HostKeys []*crypto.PrivateKey
	// Metadata
	Metadata map[string]string
	// Middlewares
	Middlewares []Handler
	// change observer
	observer chan struct{}
}

// newOptions creates new option
func newOptions(opts ...Option) *Options {
	o := &Options{
		ServerID:               uuid.MustV4(),
		NoClientAuth:           false,
		PasswordAuthentication: false,
		RSAAuthentication:      false,
		ListenAddr:             ":0",
		Protocol:               2,
		HostKeys:               make([]*crypto.PrivateKey, 0),
		Metadata:               make(map[string]string, 0),
		observer:               make(chan struct{}),
	}
	for _, opt := range opts {
		opt(o)
	}
	if len(o.HostKeys) == 0 {
		k, err := crypto.Generate()
		if err != nil {
			log.Fatal(err)
		}
		o.AddHostKey(k)
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
		o.AddHostKey(k)
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
	go v.notify()
	return v
}

// SetPasswordAuthentication to enable or disable password authentication
func (v *Options) SetPasswordAuthentication(enable bool) *Options {
	v.PasswordAuthentication = enable
	go v.notify()
	return v
}

// SetRSAAuthentication to enable or disable rsa key authentication
func (v *Options) SetRSAAuthentication(enable bool) *Options {
	v.RSAAuthentication = enable
	go v.notify()
	return v
}

// SetListenAddr to set listen port
func (v *Options) SetListenAddr(addr string) *Options {
	if len(addr) > 0 {
		v.ListenAddr = addr
		go v.notify()
	}
	return v
}

// SetProtocol to set secure shell protocol
func (v *Options) SetProtocol(protocol int) *Options {
	if protocol > 0 {
		v.Protocol = protocol
		go v.notify()
	}
	return v
}

// AddHostKey to set host private key
func (v *Options) AddHostKey(k *crypto.PrivateKey) *Options {
	v.HostKeys = append(v.HostKeys, k)
	go v.notify()
	return v
}

// AddMiddleware to add auth middleware
func (v *Options) AddMiddleware(fs ...Handler) *Options {
	for _, f := range fs {
		v.Middlewares = append(v.Middlewares, f)
		go v.notify()
	}
	return v
}

// SetServerID to set server reference id
func (v *Options) SetServerID(s string) *Options {
	if len(s) > 0 {
		v.ServerID = s
		go v.notify()
	}
	return v
}

// GetServerID to return server id
func (v *Options) GetServerID() string {
	return v.ServerID
}

// GetClientAuth to return client auth settings
func (v *Options) GetClientAuth() bool {
	return !v.NoClientAuth
}

// GetPasswordAuthentication to return password authentication setting
func (v *Options) GetPasswordAuthentication() bool {
	return v.PasswordAuthentication
}

// GetRSAAuthentication to return rsa authentication setting
func (v *Options) GetRSAAuthentication() bool {
	return v.RSAAuthentication
}

// GetListenAddr to return listen address
func (v *Options) GetListenAddr() string {
	return v.ListenAddr
}

// GetProtocol to return protocol version
func (v *Options) GetProtocol() int {
	return v.Protocol
}

// GetHostKeys to return host private key
func (v *Options) GetHostKeys() []*crypto.PrivateKey {
	return v.HostKeys
}

// GetMetadata to return metadata data
func (v *Options) GetMetadata(k string) string {
	v.RLock()
	defer v.RUnlock()
	return v.Metadata[k]
}

// GetMiddlewares to return middlewares
func (v *Options) GetMiddlewares() []Handler {
	return v.Middlewares
}

func (v *Options) notify() {
	v.observer <- struct{}{}
}

// Changes observer
func (v *Options) Changes() <-chan struct{} {
	return v.observer
}
