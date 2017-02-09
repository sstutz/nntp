package nntp

import (
	"reflect"
	"testing"
)

var (
	addr string = "127.0.0.1"
	port int    = 0
)

func TestNewSocket(t *testing.T) {
	want := &Socket{
		Address: addr,
		Port:    port,
		ssl:     false,
	}
	if socket := NewSocket(addr, port); !reflect.DeepEqual(socket, want) {
		t.Errorf("NewSocket() = %+v, expected %+v", socket, want)
	}
}

func TestNewSSLSocket(t *testing.T) {
	want := &Socket{
		Address: addr,
		Port:    port,
		ssl:     true,
	}
	if socket := NewSSLSocket(addr, port); !reflect.DeepEqual(socket, want) {
		t.Errorf("NewSSLSocket() = %+v, expected %+v", socket, want)
	}
}

func TestSocketToString(t *testing.T) {
	want := "127.0.0.1:0"
	socket := NewSocket(addr, port)
	str := socket.String()
	if str != want {
		t.Errorf("String() = %+v, expected %+v", str, want)
	}
}
