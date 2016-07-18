package nntp

import (
	"bufio"
	"io"
	"net/textproto"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
)

// NNTP DATE command
// yyyymmddhhmmss Server date and time
const DateFormat = "20060102150405"

type Client struct {
	connection *textproto.Conn
	Banner     string
	compress   bool
}

type Credentials struct {
	Username, Password string
}

type Response struct {
	Code int
	Line string
	Data io.Reader
}

func (c *Client) Article(messageId string) (*Article, error) {
	response, err := c.multilineCommand("ARTICLE "+messageId, 220)
	if err != nil {
		return nil, err
	}

	be := bufio.NewReader(response.Data)
	tp := textproto.NewReader(be)
	header, err := tp.ReadMIMEHeader()
	if err != nil && err != io.EOF {
		return nil, err
	}

	return &Article{
		Header: header,
		Body:   be,
		Bytes:  1,
		Lines:  2,
	}, nil
}

// Authenticates against an NNTP server per RFC 4643
func (c *Client) Authenticate(cred *Credentials) (err error) {
	log.Infof("Authenticate User %s", cred.Username)
	if _, err = c.command("AUTHINFO USER "+cred.Username, 381); err != nil {
		return
	}
	if _, err = c.command("AUTHINFO PASS "+cred.Password, 281); err != nil {
		return
	}

	return
}

func (c *Client) Body(messageId string) (io.Reader, error) {
	response, err := c.multilineCommand("BODY"+messageId, 222)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

// Returns  a short summary of the commands that are understood by this implementation of the server
func (c *Client) Capabilities() ([]string, error) {
	response, err := c.multilineCommand("CAPABILITIES", 101)
	if err != nil {
		return nil, err
	}
	capabilities := []string{}
	scanner := bufio.NewScanner(response.Data)
	for scanner.Scan() {
		capabilities = append(capabilities, scanner.Text())
	}
	if scanner.Err() != nil {
		return nil, scanner.Err()
	}

	return capabilities, nil
}

// Terminates current sessions
func (c *Client) Close() error {
	if _, err := c.command("QUIT", 206); err != nil {
		return err
	}

	return c.connection.Close()
}

func (c *Client) Date() (t time.Time, err error) {
	var response *Response
	if response, err = c.command("DATE", 111); err != nil {
		return
	}
	if t, err = time.Parse(DateFormat, response.Line); err != nil {
		return
	}

	return
}

// Selects a newsgroup as the currently selected newsgroup
func (c *Client) Group(name string) (*Group, error) {
	response, err := c.command("GROUP "+name, 211)
	if err != nil {
		return nil, err
	}
	group, err := groupFromLine(response.Line)
	log.Infof("Group selected: %s", group.Name)
	return group, err
}

type GroupListing struct {
	Number, Low, High int
	Group             string
	Range             []int
}

func (c *Client) ListGroup(group string) (*GroupListing, error) {
	response, err := c.multilineCommand("LISTGROUP "+group, 211)
	if err != nil {
		return nil, err
	}

	params := strings.SplitN(response.Line, " ", 4)
	n, _ := strconv.Atoi(params[0])
	l, _ := strconv.Atoi(params[1])
	h, _ := strconv.Atoi(params[2])

	r := []int{}
	scanner := bufio.NewScanner(response.Data)
	for scanner.Scan() {
		article, _ := strconv.Atoi(scanner.Text())
		r = append(r, article)
	}

	return &GroupListing{
		Number: n,
		Low:    l,
		High:   h,
		Group:  params[3],
		Range:  r,
	}, scanner.Err()
}

//
func (c *Client) Head(message string) (header textproto.MIMEHeader, err error) {
	var response *Response
	if response, err = c.multilineCommand("HEAD "+message, 221); err != nil {
		return
	}
	// MIME Headers cannot be extracted from a dot encoded block / text
	// This line returns a new Reader that satisfie reads using the decoded text
	tp := textproto.NewReader(bufio.NewReader(response.Data))
	if header, err = tp.ReadMIMEHeader(); err != nil && err != io.EOF {
		return
	}

	return header, nil
}

// Returns a short summary of the commands that are understood by this implementation of the server
func (c *Client) Help() ([]string, error) {
	response, err := c.multilineCommand("HELP", 100)
	if err != nil {
		return nil, err
	}
	help := []string{}
	scanner := bufio.NewScanner(response.Data)
	for scanner.Scan() {
		help = append(help, scanner.Text())
	}
	if scanner.Err() != nil {
		return nil, scanner.Err()
	}

	return help, nil
}

// Selects the previous article
func (c *Client) Last() (string, string, error) {
	response, err := c.command("LAST", 223)
	if err != nil {
		return "", "", err
	}
	params := strings.SplitN(response.Line, " ", 3)
	if len(params) < 2 {
		return "", "", textproto.ProtocolError("Unexpected end of LAST")
	}

	return params[0], params[1], nil
}

// Returns a list of all available groups
func (c *Client) List() ([]Group, error) {
	response, err := c.multilineCommand("LIST", 215)
	if err != nil {
		return nil, err
	}
	groups := []Group{}
	scanner := bufio.NewScanner(response.Data)
	for scanner.Scan() {
		group, err := groupFromListLine(scanner.Text())
		if err != nil {
			return nil, err
		}
		groups = append(groups, *group)
	}
	if scanner.Err() != nil {
		return nil, scanner.Err()
	}

	return groups, nil
}

// Returns a list of new groups since N
func (c *Client) NewGroups() error {
	return nil
	//@todo(sstutz): Implement NEWGROUPS
}

func (c *Client) NewNews() error {
	return nil
	//@todo(sstutz): Implement NEWNEWS
}

// Selects the next article
func (c *Client) Next() (string, string, error) {
	response, err := c.command("NEXT", 223)
	if err != nil {
		return "", "", nil
	}
	params := strings.SplitN(response.Line, " ", 3)
	if len(params) < 2 {
		return "", "", textproto.ProtocolError("Unexpected end of NEXT")
	}

	return params[0], params[1], nil
}

// Behaves like Article except that, if the article exists, it is not presented to the client
func (c *Client) Stat(id string) (string, string, error) {
	response, err := c.command("STAT "+id, 223)
	if err != nil {
		return "", "", err
	}
	params := strings.SplitN(response.Line, " ", 3)
	if len(params) < 2 {
		return "", "", textproto.ProtocolError("Unexpected end of STAT")
	}

	return params[0], params[1], nil
}

func (c *Client) Xover() error {
	//@todo(sstutz): Implement XOVER https://tools.ietf.org/html/rfc2980
	return nil
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
	if err = c.connection.PrintfLine(cmd); err != nil {
		return
	}
	if rc, line, err = c.connection.ReadCodeLine(expected); err != nil {
		return
	}
	return
}
