package client

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"net/textproto"
	"strconv"
	"strings"

	"github.com/davidklassen/wow-pow/pkg/challenge"
)

type Client struct {
	addr string
	conn net.Conn
}

func New(addr string) *Client {
	return &Client{addr: addr}
}

func (c *Client) Connect() error {
	var err error
	c.conn, err = net.Dial("tcp", c.addr)
	if err != nil {
		return fmt.Errorf("failed to dial tcp: %w", err)
	}
	return nil
}

func (c *Client) Quote() (string, error) {
	if _, err := c.conn.Write([]byte("get\n")); err != nil {
		return "", fmt.Errorf("failed to write command: %w", err)
	}

	reader := textproto.NewReader(bufio.NewReader(c.conn))
	data, err := reader.ReadLine()
	if err != nil {
		return "", fmt.Errorf("failed to read challenge: %w", err)
	}

	parts := strings.Split(data, ":")
	if len(parts) != 2 {
		return "", errors.New("invalid challenge format")
	}

	bits, err := strconv.Atoi(parts[0])
	if err != nil {
		return "", fmt.Errorf("failed to parse bits: %w", err)
	}

	if _, err = fmt.Fprintf(c.conn, "%s\n", challenge.Solve(parts[1], bits)); err != nil {
		return "", fmt.Errorf("failed to write solution: %w", err)
	}

	var quote strings.Builder
	for {
		line, err := reader.ReadLine()
		if err != nil {
			return "", fmt.Errorf("failed to read quote line: %w", err)
		}
		if line == "" {
			break
		}
		quote.WriteString(line)
		quote.WriteByte('\n')
	}

	return quote.String(), nil
}
