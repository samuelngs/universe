package server

import (
	"log"

	"golang.org/x/crypto/ssh"
)

// DebugHandler logging handler
type DebugHandler func(md ssh.ConnMetadata, method string, err error)

// Configs struct
type Configs struct {
	conf *ssh.ServerConfig
	opts *Options
	logs DebugHandler
}

func newConfigs(opts *Options) *Configs {
	c := new(Configs)
	c.opts = opts
	c.conf = &ssh.ServerConfig{
		PasswordCallback:  c.PasswordCallback,
		PublicKeyCallback: c.PublicKeyCallback,
		AuthLogCallback:   c.AuthLogCallback,
	}
	for _, key := range opts.GetHostKeys() {
		signer, err := key.Signer()
		if err != nil {
			log.Fatal(err)
		}
		c.conf.AddHostKey(signer)
	}
	go c.sync(opts)
	return c
}

func (v *Configs) sync(opts *Options) {
	for {
		select {
		case <-opts.Changes():
			v.conf = &ssh.ServerConfig{
				PasswordCallback:  v.PasswordCallback,
				PublicKeyCallback: v.PublicKeyCallback,
				AuthLogCallback:   v.AuthLogCallback,
			}
			for _, key := range opts.GetHostKeys() {
				signer, err := key.Signer()
				if err != nil {
					log.Fatal(err)
				}
				v.conf.AddHostKey(signer)
			}
		}
	}
}

// PasswordCallback func
func (v *Configs) PasswordCallback(md ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
	switch {
	case v.opts.NoClientAuth:
		return nil, nil
	case v.opts.PasswordAuthentication:
		c := &Context{
			typ:   AuthenticationPassword,
			raddr: md.RemoteAddr(),
			laddr: md.LocalAddr(),
		}
		for _, handler := range v.opts.GetMiddlewares() {
			if err := handler(c); err != nil {
				return nil, err
			}
		}
		return nil, nil
	default:
		return nil, ErrUnauthentized
	}
}

// PublicKeyCallback func
func (v *Configs) PublicKeyCallback(md ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
	switch {
	case v.opts.NoClientAuth:
		return nil, nil
	case v.opts.RSAAuthentication:
		c := &Context{
			typ:   AuthenticationPublicKey,
			raddr: md.RemoteAddr(),
			laddr: md.LocalAddr(),
		}
		for _, handler := range v.opts.GetMiddlewares() {
			if err := handler(c); err != nil {
				return nil, err
			}
		}
		return nil, nil
	default:
		return nil, ErrUnauthentized
	}
}

// AuthLogCallback func
func (v *Configs) AuthLogCallback(md ssh.ConnMetadata, method string, err error) {
	if v.logs != nil {
		v.logs(md, method, err)
	}
}
