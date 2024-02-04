package server

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/textproto"
	"sync"
	"time"

	"github.com/davidklassen/wow-pow/pkg/challenge"
)

const (
	// challenge data length.
	dataLen = 12
)

// Handler interface defines the contract for handling different commands received by the server.
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

// serve handles individual client connections.
// It reads the client's request, sends a PoW challenge, verifies the solution,
// and then processes the command if the verification is successful.
func (s *Server) serve(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			slog.Error("failed to close connection", slog.String("error", err.Error()))
		}
		s.wg.Done()
	}()

	reader := textproto.NewReader(bufio.NewReader(conn))
	for {
		start := time.Now()

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

		data := challenge.Generate(dataLen, s.difficulty)
		if _, err = fmt.Fprintf(conn, "%s\n", data); err != nil {
			slog.Error("failed to send challenge", slog.String("error", err.Error()))
			return
		}

		solution, err := reader.ReadLine()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				slog.Error("failed to read solution", slog.String("error", err.Error()))
			}
			return
		}

		if err = challenge.Verify(data, solution); err != nil {
			slog.Warn("failed to verify solution", slog.String("error", err.Error()))
			return
		}

		if err = s.handler.Handle(cmd, conn); err != nil {
			slog.Error("failed to handle request", slog.String("error", err.Error()))
			return
		}

		slog.Info("request handled",
			slog.String("remote_address", conn.RemoteAddr().String()),
			slog.Duration("duration", time.Since(start)),
		)

		// FIXME: should check for stopped state
	}
}

func (s *Server) Start() error {
	var err error
	s.listener, err = net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("failed to create listener: %w", err)
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
	return nil
}

func (s *Server) Stop() error {
	if err := s.listener.Close(); err != nil {
		return fmt.Errorf("failed to close listener: %w", err)
	}
	slog.Info("stopped listening", slog.String("address", s.addr))
	s.wg.Wait()
	slog.Info("server stopped")
	return nil
}
