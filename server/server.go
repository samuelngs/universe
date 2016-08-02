package server

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/crypto/ssh"
)

// Server daemon for Secure Shell
type Server interface {
	Use(...Handler)
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
	ser.config = newConfigs(ser.option)
	// ser.config = &ssh.ServerConfig{
	// 	AuthLogCallback: func(md ssh.ConnMetadata, method string, err error) {
	// 		switch {
	// 		case err != nil:
	// 			ser.logger <- &trace{
	// 				topic:   TraceAuthentication,
	// 				message: fmt.Sprintf("Reject connection from %s@%s, %s (%s)", md.User(), md.LocalAddr(), md.RemoteAddr(), md.ClientVersion()),
	// 				err:     err,
	// 			}
	// 		default:
	// 			ser.logger <- &trace{
	// 				topic:   TraceAuthentication,
	// 				message: fmt.Sprintf("Accept connection from %s@%s, %s (%s)", md.User(), md.LocalAddr(), md.RemoteAddr(), md.ClientVersion()),
	// 			}
	// 		}
	// 	},
	// }
	return ser
}

// internal server
type server struct {
	option  *Options
	config  *Configs
	events  chan Event
	logger  chan Log
	started bool
}

func (v *server) observe(listener net.Listener) {
	v.started = true
	for {
		tcpconn, err := listener.Accept()
		if err != nil {
			continue
		}
		sshconn, chans, reqs, err := ssh.NewServerConn(tcpconn, v.config.conf)
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

func (v *server) Use(fs ...Handler) {
	v.option.AddMiddleware(fs...)
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
