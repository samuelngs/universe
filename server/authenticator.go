package server

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

type authenticator struct {
	option *Options
	logger chan Log
}

func (v *authenticator) Password(md ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
	v.logger <- &trace{
		topic:   TracePasswordAuthentication,
		message: fmt.Sprintf("Password authentication from %s@%s, %s (%s) [%s]", md.User(), md.LocalAddr(), md.RemoteAddr(), md.ClientVersion(), string(pass[:])),
	}
	return nil, nil
}

func (v *authenticator) PublicKey(md ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
	v.logger <- &trace{
		topic:   TraceRSAAuthentication,
		message: fmt.Sprintf("RSA authentication from %s@%s, %s (%s) [%s]", md.User(), md.LocalAddr(), md.RemoteAddr(), md.ClientVersion(), string(key.Marshal()[:])),
	}
	return nil, nil
}
