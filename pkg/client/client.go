package client

import (
	"bufio"
	"fmt"
	"net"
	"net/textproto"
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

func (c *Client) connect() error {
	var err error
	c.conn, err = net.Dial("tcp", c.addr)
	if err != nil {
		return fmt.Errorf("failed to dial tcp: %w", err)
	}
	return nil
}

// Quote retrieves a quote from the server.
func (c *Client) Quote() (string, error) {
	if c.conn == nil {
		if err := c.connect(); err != nil {
			return "", fmt.Errorf("failed to establish connection: %w", err)
		}
	}

	if _, err := c.conn.Write([]byte("get\n")); err != nil {
		// handle disconnection
		_ = c.conn.Close()
		c.conn = nil
		return "", fmt.Errorf("failed to write command: %w", err)
	}

	reader := textproto.NewReader(bufio.NewReader(c.conn))
	data, err := reader.ReadLine()
	if err != nil {
		return "", fmt.Errorf("failed to read challenge: %w", err)
	}

	res, err := challenge.Solve(data)
	if err != nil {
		return "", fmt.Errorf("failed to solve challenge: %w", err)
	}

	if _, err = fmt.Fprintf(c.conn, "%s\n", res); err != nil {
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
