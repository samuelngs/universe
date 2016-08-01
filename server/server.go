package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/samuelngs/universe/pkg/crypto"

	"golang.org/x/crypto/ssh"
)

// type Handler

// Server daemon for Secure Shell
type Server interface {
	Run() error
	Stop() error
	Started() bool
	Option() *Options
	Subscribe() <-chan Event
	Logging() <-chan Log
}

// New create secure shell server
func New(opts ...Option) Server {
	ser := new(server)
	ser.events = make(chan Event)
	ser.logger = make(chan Log)
	ser.option = newOptions(opts...)
	ser.authenticator = &authenticator{
		option: ser.option,
		logger: ser.logger,
	}
	ser.config = &ssh.ServerConfig{
		PasswordCallback: func(md ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			switch {
			case ser.option.NoClientAuth:
				return nil, nil
			case ser.option.PasswordAuthentication:
				return ser.authenticator.Password(md, pass)
			default:
				return nil, ErrUnauthentized
			}
		},
		PublicKeyCallback: func(md ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
			switch {
			case ser.option.NoClientAuth:
				return nil, nil
			case ser.option.RSAAuthentication:
				return ser.authenticator.PublicKey(md, key)
			default:
				return nil, ErrUnauthentized
			}
		},
		AuthLogCallback: func(md ssh.ConnMetadata, method string, err error) {
			switch {
			case err != nil:
				ser.logger <- &trace{
					topic:   TraceAuthentication,
					message: fmt.Sprintf("Reject connection from %s@%s, %s (%s)", md.User(), md.LocalAddr(), md.RemoteAddr(), md.ClientVersion()),
					err:     err,
				}
			default:
				ser.logger <- &trace{
					topic:   TraceAuthentication,
					message: fmt.Sprintf("Accept connection from %s@%s, %s (%s)", md.User(), md.LocalAddr(), md.RemoteAddr(), md.ClientVersion()),
				}
			}
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
	option        *Options
	config        *ssh.ServerConfig
	events        chan Event
	logger        chan Log
	started       bool
	authenticator *authenticator
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
			v.logger <- &trace{
				topic:   TraceHandshake,
				message: fmt.Sprintf("Failed to handshake %v", tcpconn.RemoteAddr()),
				err:     err,
			}
			continue
		}
		v.logger <- &trace{
			topic:   TraceConnect,
			message: fmt.Sprintf("New connection from %s (%s)", sshconn.RemoteAddr(), sshconn.ClientVersion()),
		}
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
	if typ := channel.ChannelType(); typ != "session" {
		s := fmt.Sprintf("Unknown channel type: %s", typ)
		channel.Reject(ssh.UnknownChannelType, s)
		v.logger <- &trace{
			topic:   TraceChannel,
			message: s,
		}
		return
	}
	connection, requests, err := channel.Accept()
	if err != nil {
		v.logger <- &trace{
			topic:   TraceChannel,
			message: "Could not accept channel",
			err:     err,
		}
	}
	// defer connection.Close()
	if connection != nil {
	}
	go v.process(requests)
}

func (v *server) process(reqs <-chan *ssh.Request) {
	for req := range reqs {
		switch {
		case req.Type == "shell":
		}
		fmt.Println(req)
	}
}

func (v *server) Run() error {
	v.events <- &event{
		topic:   EventServerStart,
		message: "Starting server",
	}
	listener, err := net.Listen("tcp", v.option.ListenAddr)
	if err != nil {
		return err
	}
	v.events <- &event{
		topic:   EventServerStarted,
		message: fmt.Sprintf("Listening on %s", listener.Addr().String()),
	}
	defer listener.Close()
	go v.observe(listener)
	ch := make(chan os.Signal, 1)
	signal.Notify(
		ch,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGKILL,
	)
	v.events <- &event{
		topic:   EventReceiveSignal,
		message: fmt.Sprintf("Received signal %s", <-ch),
	}
	return v.Stop()
}

func (v *server) Stop() error {
	v.events <- &event{
		topic:   EventServerStop,
		message: "Stopping server",
	}
	v.started = false
	v.events <- &event{
		topic: EventServerStopped,
	}
	defer close(v.events)
	defer close(v.logger)
	return nil
}

func (v *server) Started() bool {
	return v.started
}

func (v *server) Option() *Options {
	return v.option
}

func (v *server) Subscribe() <-chan Event {
	return v.events
}

func (v *server) Logging() <-chan Log {
	return v.logger
}
