package nntp

import "io"

type Response struct {
	Code int
	Line string
	Data io.Reader
}

// sends a low level command to nntp and checks the response code against the expected one.
func (c *Client) command(cmd string, expected int) (*Response, error) {
	rc, line, err := c.sendCommand(cmd, expected)
	if err != nil {
		return nil, err
	}
	return &Response{
		Code: rc,
		Line: line,
	}, nil
}

// sends a low level command to nntp and checks the response code against the expected one.
// Also handles Multi-line Data Block reponse
func (c *Client) multilineCommand(cmd string, expected int) (*Response, error) {
	rc, line, err := c.sendCommand(cmd, expected)
	if err != nil {
		return nil, err
	}

	return &Response{
		Code: rc,
		Line: line,
		Data: c.connection.DotReader(),
	}, nil
}

// Helper to send the actual command and checks the reponse code against the expexted one.
func (c *Client) sendCommand(cmd string, expected int) (rc int, line string, err error) {
	println("test in sendcommand")
	if err = c.connection.PrintfLine(cmd); err != nil {
		return
	}
	if rc, line, err = c.connection.ReadCodeLine(expected); err != nil {
		return
	}
	return
}
