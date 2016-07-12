package nntp

import (
	"crypto/tls"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"net/textproto"
)

func New(s *Socket) (client *Client, err error) {
	c := new(textproto.Conn)
	if s.ssl {
		if c, err = dialTLS(s.String()); err != nil {
			return nil, err
		}
	} else {
		if c, err = dial(s.String()); err != nil {
			return nil, err
		}
	}
	log.Infof("Connected to NNTP")
	_, msg, err := c.ReadCodeLine(200)
	if err != nil {
		return nil, err
	}

	return &Client{
		connection: c,
		Banner:     msg,
	}, nil
}

func dial(addr string) (*textproto.Conn, error) {
	c, err := textproto.Dial("tcp", addr)

	if err != nil {
		return nil, err
	}

	return c, nil

}

func dialTLS(addr string) (*textproto.Conn, error) {
	fmt.Println("[DialTLS] Connecting over TLS")
	c, err := tls.Dial("tcp", addr, nil)

	if err != nil {
		return nil, err
	}

	return textproto.NewConn(c), nil
}
