package nntp

import (
	"net"
	"strconv"
)

type Socket struct {
	Address string
	Port    int
	ssl     bool
}

func (s Socket) String() string {
	return net.JoinHostPort(s.Address, strconv.Itoa(s.Port))
}

func NewSocket(address string, port int) *Socket {
	return &Socket{
		Address: address,
		Port:    port,
		ssl:     false,
	}
}

func NewSSLSocket(address string, port int) *Socket {
	return &Socket{
		Address: address,
		Port:    port,
		ssl:     true,
	}
}
