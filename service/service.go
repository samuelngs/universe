package service

import "github.com/samuelngs/universe/sshd"

// Service for Secure Shell
type Service interface {
	Run() error
	Stop() error
	Server() sshd.Server
	Client() client.Client
}
