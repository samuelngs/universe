package server

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/samuelngs/universe/pkg/crypto"

	"golang.org/x/crypto/ssh"
)

// Server daemon for Secure Shell
type Server interface {
	Run() error
	Stop() error
	Started() bool
	Option() *Option
	Verify() error
	Subscribe() chan<- Event
}

// Event interface for secure shell server
type Event interface {
	Topic() string
	Message() interface{}
}

// New create secure shell server
func New(opts ...Option) Server {
	opt := &Option{
		NoClientAuth:           false,
		PasswordAuthentication: false,
		RSAAuthentication:      false,
		Addr:                   ":0",
		Protocol:               2,
		HostKey:                nil,
	}
	for _, o := range opts {
		opt = &o
		break
	}
	ser := new(server)
	ser.events = make(chan Event)
	ser.option = opt
	ser.config = &ssh.ServerConfig{
		PasswordCallback: func(md ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			switch {
			case ser.option.NoClientAuth:
				return nil, nil
			case ser.option.PasswordAuthentication == false:
				return nil, ErrUnauthentized
			case ser.option.PasswordAuthentication:
			}
			return nil, ErrUnauthentized
		},
		PublicKeyCallback: func(md ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
			switch {
			case ser.option.NoClientAuth:
				return nil, nil
			case ser.option.RSAAuthentication == false:
				return nil, ErrUnauthentized
			case ser.option.RSAAuthentication:
			}
			return nil, ErrUnauthentized
		},
	}
	if ser.option.HostKey == nil {
		k, err := crypto.Generate()
		if err != nil {
			log.Fatal(err)
		}
		ser.option.HostKey = k
	}
	if signer, err := ser.option.HostKey.Signer(); err != nil {
		log.Fatal(err)
	} else if err == nil {
		ser.config.AddHostKey(signer)
	}
	return ser
}

// internal server
type server struct {
	option  *Option
	config  *ssh.ServerConfig
	events  chan<- Event
	started bool
}

func (v *server) observe(listener net.Listener) {
	v.started = true
	for {
		tcpconn, err := listener.Accept()
		if err != nil {
			continue
		}
		sshconn, chans, reqs, err := ssh.NewServerConn(tcpconn, v.config)
		if err != nil {
			log.Printf("Failed to handshake (%s)", err)
			continue
		}
		log.Printf("New connection from %s (%s)", sshconn.RemoteAddr(), sshconn.ClientVersion())
		go ssh.DiscardRequests(reqs)
		go v.receiver(chans)
	}
}

func (v *server) receiver(chans <-chan ssh.NewChannel) {
	for channel := range chans {
		go v.handle(channel)
	}
}

func (v *server) handle(channel ssh.NewChannel) {
	log.Printf("handle channel")
}

func (v *server) Run() error {
	log.Printf("Starting server")
	listener, err := net.Listen("tcp", v.option.Addr)
	if err != nil {
		return err
	}
	defer listener.Close()
	log.Printf("Listening on %s", listener.Addr().String())
	go v.observe(listener)
	ch := make(chan os.Signal, 1)
	signal.Notify(
		ch,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGKILL,
	)
	log.Printf("Received signal %s", <-ch)
	return v.Stop()
}

func (v *server) Stop() error {
	log.Printf("Stopping server")
	v.started = false
	return nil
}

func (v *server) Started() bool {
	return v.started
}

func (v *server) Option() *Option {
	return v.option
}

func (v *server) Subscribe() chan<- Event {
	return v.events
}

func (v *server) Verify() error {
	return nil
}
