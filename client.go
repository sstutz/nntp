package nntp

import (
	"bufio"
	"fmt"
	"io"
	"net/textproto"
	"strings"

	log "github.com/Sirupsen/logrus"
)

type Client struct {
	connection *textproto.Conn
	Banner     string
	compress   bool
}

type Credentials struct {
	Username, Password string
}

// Authenticates against an NNTP server per RFC 4643
func (c *Client) Authenticate(cred *Credentials) (err error) {
	log.Infof("Authenticate User %s", cred.Username)
	if _, _, err = c.command(fmt.Sprintf("AUTHINFO USER %s", cred.Username), 381); err != nil {
		return
	}

	if _, _, err = c.command(fmt.Sprintf("AUTHINFO PASS %s", cred.Password), 281); err != nil {
		return
	}

	return
}

// Returns  a short summary of the commands that are understood by this implementation of the server
func (c *Client) Capabilities() (lines []string, err error) {
	if _, lines, err = c.multilineCommand("CAPABILITIES", 101); err != nil {
		return
	}
	return
}

// Terminates current sessions
func (c *Client) Close() (err error) {
	if _, _, err = c.command("QUIT", 206); err != nil {
		return err
	}

	c.connection.Close()
	return
}

// Selects a newsgroup as the currently selected newsgroup
func (c *Client) Group(name string) (group *Group, err error) {
	_, line, err := c.command(fmt.Sprintf("GROUP %s", name), 211)
	if err != nil {
		return nil, err
	}
	group, err = groupFromLine(line)

	log.Infof("Group selected: %s", group.Name)

	return group, err
}

func (c *Client) Head(message string) (header textproto.MIMEHeader, err error) {
	if err = c.connection.PrintfLine("HEAD " + message); err != nil {
		return
	}

	if _, _, err = c.connection.ReadCodeLine(221); err != nil {
		return
	}

	// MIME Headers cannot be extracted from a dot encoded block / text
	// This line returns a new Reader that satisfie reads using the decoded text
	tp := textproto.NewReader(bufio.NewReader(c.connection.DotReader()))
	if header, err = tp.ReadMIMEHeader(); err != nil && err != io.EOF {
		return
	}

	return header, nil
}

// Returns a short summary of the commands that are understood by this implementation of the server
func (c *Client) Help() (lines []string, err error) {
	if _, lines, err = c.multilineCommand("HELP", 100); err != nil {
		return
	}
	return
}

// Selects the previous article
func (c *Client) Last() (string, string, error) {
	_, line, err := c.command("LAST", 223)
	if err != nil {
		return "", "", nil
	}

	params := strings.SplitN(line, " ", 3)
	if len(params) < 2 {
		return "", "", textproto.ProtocolError("Unexpected end of LAST")
	}

	return params[0], params[1], nil
}

func (c *Client) List() ([]Group, error) {
	_, lines, err := c.multilineCommand("LIST", 215)
	if err != nil {
		return nil, err
	}

	groups := make([]Group, 0, len(lines))
	for _, line := range lines[1:] {
		group, err := groupFromListLine(line)
		if err != nil {
			return nil, err
		}
		groups = append(groups, *group)
	}

	return groups, nil
}

// Selects the next article
func (c *Client) Next() (string, string, error) {
	_, line, err := c.command("NEXT", 223)
	if err != nil {
		return "", "", nil
	}

	params := strings.SplitN(line, " ", 3)
	if len(params) < 2 {
		return "", "", textproto.ProtocolError("Unexpected end of NEXT")
	}

	return params[0], params[1], nil
}

// Behaves like Article except that, if the article exists, it is not presented to the client
func (c *Client) Stat(id string) (string, string, error) {
	_, line, err := c.command(fmt.Sprintf("STAT %s", id), 223)
	if err != nil {
		return "", "", err
	}

	params := strings.SplitN(line, " ", 3)
	if len(params) < 2 {
		return "", "", textproto.ProtocolError("Unexpected end of STAT")
	}

	return params[0], params[1], nil
}

// sends a low level command to nntp and checks the response code against the expected one.
func (c *Client) command(cmd string, expected int) (rc int, line string, err error) {
	if err = c.connection.PrintfLine(cmd); err != nil {
		return 0, "", err
	}

	if rc, line, err = c.connection.ReadCodeLine(expected); err != nil {
		return rc, "", err
	}

	return
}

// sends a low level command to nntp and checks the response code against the expected one.
// Also handles Multi-line Data Block reponse
func (c *Client) multilineCommand(cmd string, expected int) (rc int, lines []string, err error) {
	if err = c.connection.PrintfLine(cmd); err != nil {
		return 0, nil, err
	}

	var line string
	if rc, line, err = c.connection.ReadCodeLine(expected); err != nil {
		return rc, nil, err
	}

	lines, err = c.connection.ReadDotLines()
	if err != nil {
		return rc, nil, err
	}

	l := []string{line}
	return rc, append(l, lines...), nil

}
