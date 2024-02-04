package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/davidklassen/wow-pow/pkg/server"
)

var (
	addr       = flag.String("addr", ":1111", "application network address")
	dbFile     = flag.String("db", "db.txt", "database file")
	difficulty = flag.Int("difficulty", 4, "challenge difficulty")
	timeout    = flag.Duration("timeout", time.Second*3, "connection idle timeout")
)

// quoteHandler manages the retrieval of quotes from the database.
// It uses a round-robin method to select a quote from the database.
type quoteHandler struct {
	db     []string
	nextID atomic.Uint32
}

func (h *quoteHandler) quote() string {
	if len(h.db) == 0 {
		return "hello, world\n--Brian Kernighan, Programming in C: A Tutorial\n"
	}
	return h.db[int(h.nextID.Add(1))%len(h.db)]
}

func (h *quoteHandler) Handle(cmd string, w io.Writer) error {
	if cmd != "get" {
		return errors.New("bad request")
	}

	if _, err := fmt.Fprintf(w, "%s\n", h.quote()); err != nil {
		return errors.New("failed to send quote")
	}

	return nil
}

func readDB(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = file.Close(); err != nil {
			slog.Error("failed to close db file", slog.String("error", err.Error()))
		}
	}()

	var res []string
	scanner := bufio.NewScanner(file)
	var quote strings.Builder
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			quote.WriteString(line)
			quote.WriteByte('\n')
		} else {
			res = append(res, quote.String())
			quote.Reset()
		}
	}
	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func main() {
	flag.Parse()

	db, err := readDB(*dbFile)
	if err != nil {
		slog.Error("failed to read quotes DB", slog.String("error", err.Error()))
	}
	h := &quoteHandler{db: db}

	srv := server.New(*addr, h, *difficulty, *timeout)
	if err = srv.Start(); err != nil {
		slog.Error("failed to start server", slog.String("error", err.Error()))
		os.Exit(1)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	slog.Info("received quit signal", slog.String("signal", (<-sig).String()))

	if err = srv.Stop(); err != nil {
		slog.Error("failed to stop server", slog.String("error", err.Error()))
		os.Exit(1)
	}
	slog.Info("bye bye")
}
