package server

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/textproto"
	"os"
	"sync"
	"time"

	"github.com/davidklassen/wow-pow/pkg/challenge"
)

const (
	// challenge data length.
	dataLen = 12
)

type Handler interface {
	Handle(string, io.Writer) error
}

type Server struct {
	wg         sync.WaitGroup
	addr       string
	listener   net.Listener
	handler    Handler
	timeout    time.Duration
	difficulty int
}

func New(addr string, handler Handler, difficulty int, timeout time.Duration) *Server {
	return &Server{
		addr:       addr,
		handler:    handler,
		timeout:    timeout,
		difficulty: difficulty,
	}
}

func (s *Server) serve(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			slog.Error("failed to close connection", slog.String("error", err.Error()))
		}
		s.wg.Done()
	}()

	reader := textproto.NewReader(bufio.NewReader(conn))
	for {
		if err := conn.SetDeadline(time.Now().Add(s.timeout)); err != nil {
			slog.Error("failed to set connection deadline")
			return
		}

		cmd, err := reader.ReadLine()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				slog.Error("failed to read request", slog.String("error", err.Error()))
			}
			return
		}

		data := challenge.Generate(dataLen)
		if _, err = fmt.Fprintf(conn, "%d:%s\n", s.difficulty, data); err != nil {
			slog.Error("failed to send challenge")
			return
		}

		solution, err := reader.ReadLine()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				slog.Error("failed to read solution", slog.String("error", err.Error()))
			}
			return
		}
		if !challenge.Verify(data, solution, s.difficulty) {
			slog.Error("incorrect solution")
			return
		}

		if err = s.handler.Handle(cmd, conn); err != nil {
			slog.Error("failed to handle request", slog.String("error", err.Error()))
			return
		}
	}
}

func (s *Server) Start() {
	var err error
	s.listener, err = net.Listen("tcp", s.addr)
	if err != nil {
		slog.Error("failed to listen", slog.String("error", err.Error()))
		os.Exit(1)
	}
	slog.Info("accepting connections", slog.String("address", s.addr))

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for {
			conn, err := s.listener.Accept()
			if err != nil {
				if !errors.Is(err, net.ErrClosed) {
					slog.Error("failed to accept connection", slog.String("error", err.Error()))
				}
				return
			}
			s.wg.Add(1)
			go s.serve(conn)
		}
	}()
}

func (s *Server) Stop() {
	if err := s.listener.Close(); err != nil {
		slog.Error("failed to close listener", slog.String("error", err.Error()))
		os.Exit(1)
	}
	s.wg.Wait()
}
