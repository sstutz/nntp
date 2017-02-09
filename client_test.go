package nntp

import (
	"net/textproto"
	"testing"
)

func TestClose(t *testing.T) {
	s, conn, closer := mockClientServerWithResponse("205 connection closed")
	defer closer(s, conn)
	c := &Client{connection: textproto.NewConn(conn)}

	c.Close()
	/*
		if err != nil {
			t.Errorf("Client.Close() failed: %+v", err)
		}
	*/
}
