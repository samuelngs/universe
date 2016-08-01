package main

import (
	"flag"
	"log"
	"os"

	"github.com/samuelngs/universe/pkg/crypto"
	"github.com/samuelngs/universe/server"
)

var (
	rsa           = flag.String("rsa-key", os.Getenv("HOME")+"/.ssh/id_rsa", "path to the private key file")
	addr          = flag.String("tcp-address", "127.0.0.1:2222", "<addr>:<port> to listen on for tcp clients")
	protocol      = flag.Int("protocol", 2, "protocol version")
	noauth        = flag.Bool("disable-authentication", false, "disable authentication")
	allowpassword = flag.Bool("password-authentication", false, "allow password authentication")
	allowrsa      = flag.Bool("rsa-authentication", true, "allow rsa key authentication")
)

func main() {

	flag.Parse()

	key, err := crypto.Import(*rsa)
	if err != nil {
		log.Fatal(err)
	}

	ser := server.New(server.Option{
		NoClientAuth:           *noauth,
		PasswordAuthentication: *allowpassword,
		RSAAuthentication:      *allowrsa,
		Addr:                   *addr,
		Protocol:               *protocol,
		HostKey:                key,
	})

	ser.Run()
}
