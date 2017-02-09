package nntp

import (
	"crypto/tls"
	"net/textproto"
)

var d dialer = &netDialer{}

type dialer interface {
	dial(addr string) (*textproto.Conn, error)
	dialTLS(addr string) (*textproto.Conn, error)
}

// New creates a NNTP Client
func New(s *Socket) (client *Client, err error) {
	c := new(textproto.Conn)
	if s.ssl {
		c, err = d.dialTLS(s.String())
	} else {
		c, err = d.dial(s.String())
	}
	if err != nil {
		return nil, err
	}
	_, msg, err := c.ReadCodeLine(200)
	if err != nil {
		return nil, err
	}

	return &Client{
		connection: c,
		Banner:     msg,
	}, nil
}

type netDialer struct{}

func (n *netDialer) dial(addr string) (c *textproto.Conn, err error) {
	return textproto.Dial("tcp", addr)
}

func (n *netDialer) dialTLS(addr string) (*textproto.Conn, error) {
	c, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		return nil, err
	}
	return textproto.NewConn(c), nil
}
