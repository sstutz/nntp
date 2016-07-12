package nntp

import (
	"net/textproto"
	"strconv"
	"strings"
)

type Group struct {
	Name   string
	Number int
	High   int
	Low    int
}

func groupFromLine(line string) (*Group, error) {
	var number, high, low int
	var err error

	params := strings.SplitN(strings.TrimSpace(line), " ", 4)

	if len(params) < 4 {
		return nil, textproto.ProtocolError("unexpected end of GROUP line")
	}

	if number, err = strconv.Atoi(params[0]); err != nil {
		return nil, textproto.ProtocolError("")
	}

	if low, err = strconv.Atoi(params[1]); err != nil {
		return nil, textproto.ProtocolError("")
	}

	if high, err = strconv.Atoi(params[2]); err != nil {
		return nil, textproto.ProtocolError("")
	}

	return &Group{
		Name:   params[3],
		Number: number,
		Low:    low,
		High:   high,
	}, nil
}
