package nntp

import (
	"net"
	"net/textproto"
	"testing"
)

type mockDialer struct {
	conn net.Conn
}

func (m *mockDialer) dial(addr string) (*textproto.Conn, error) {
	return textproto.NewConn(m.conn), nil
}

func (m *mockDialer) dialTLS(addr string) (*textproto.Conn, error) {
	return textproto.NewConn(m.conn), nil
}

func mockClientServerWithResponse(r string) (net.Conn, net.Conn, func(s, c net.Conn) error) {
	s, c := net.Pipe()
	closer := func(s, c net.Conn) error {
		if s.Close() != nil || c.Close() != nil {
			println("dafuq")
		}
		return nil
	}
	go func(s net.Conn) error {
		buf := make([]byte, 1024)
		s.Read(buf)
		println("YAY" + string(buf))
		return nil
	}(s)
	return s, c, closer
}

func mockClient(tls bool, response string) (*Client, error) {
	ser, c, closer := mockClientServerWithResponse(response)
	defer closer(ser, c)
	d = &mockDialer{c}
	s := &Socket{
		Address: "0.0.0.0",
		Port:    0,
		ssl:     tls,
	}
	return New(s)
}

func newMockClient(response string) (*Client, error) {
	return mockClient(false, response)
}

func newMockTlsClient(response string) (*Client, error) {
	return mockClient(true, response)
}

func TestNew(t *testing.T) {
	response := "200 Server ready, posting allowed"
	_, err := newMockClient(response)
	_, err1 := newMockTlsClient(response)

	if err == nil {
		err = err1
	}

	if err != nil {
		t.Errorf("New(): Error: %+v", err)
	}
}
