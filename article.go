package nntp

import (
	"io"
	"net/textproto"
	"strconv"
	"strings"
)

type ArticleOverview struct {
	ArticleId  int
	Subject    string
	From       string
	Date       string
	MessageId  string
	References []string
	Bytes      int
	Lines      int
}

type Article struct {
	Header textproto.MIMEHeader
	Body   io.Reader
	// Number of bytes in the article body (used by OVER/XOVER)
	Bytes int
	// Number of lines in the article body (used by OVER/XOVER)
	Lines int
}

// MessageID provides convenient access to the article's Message ID.
func (a *Article) MessageID() string {
	return a.Header.Get("Message-Id")
}

func overviewFromLine(line string) (o *ArticleOverview, err error) {
	params := strings.SplitN(line, "\t", 9)
	var a, b, l int
	a, err = strconv.Atoi(params[0])
	b, err = strconv.Atoi(params[6])
	l, err = strconv.Atoi(params[7])
	if err != nil {
		return nil, err
	}
	o = &ArticleOverview{
		ArticleId:  a,
		Subject:    params[1],
		From:       params[2],
		Date:       params[3],
		MessageId:  params[4],
		References: strings.Split(params[5], " "),
		Bytes:      b,
		Lines:      l,
	}
	return o, nil
}
