// Package main implements the shortlink server.
//
// Note for candidates:
// This implementation of the server is not important, you should spend your
// time in linkfile.go :)
package main

import (
	"cmp"
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	serve    = flag.Bool("serve", false, "Run the web server")
	linkfile = flag.String("linkfile", "links.txt", "Filename containing shortlinks")
	addr     = flag.String("addr",
		// either ":${PORT}" or ":http"
		fmt.Sprintf(":%s", cmp.Or(os.Getenv("PORT"), "http")),
		"The address to listen on for HTTP requests.")
)

// Shortlinks is implemented by types that can convert short links to long links.
type Shortlinks interface {
	// Long returns the long URL corresponding to the short name.
	// If no URL exists for short, Long returns ShortlinkNotFoundError.
	Long(ctx context.Context, short string) (*url.URL, error)
}

func main() {
	// Parse command-line flags.
	flag.Parse()

	// Set up cancellation.
	ctx := context.Background()
	ctx, ccl := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer ccl()

	// Validate the flags.
	if !fileExists(*linkfile) {
		_, _ = fmt.Fprintln(os.Stderr, "-linkfile references a file that does not exist")
		os.Exit(1)
	}

	if *serve {
		// Creates an HTTP handler that uses the LinkFile type to find long links.
		http.Handle("GET /{short}", &ShortlinkHandler{
			&LinkFile{
				Filename: *linkfile,
			},
		})

		// create an HTTP server
		server := &http.Server{
			Addr: *addr,
		}

		if err := runServer(ctx, server); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}

	// No option selected, print usage and exit.
	flag.Usage()
	os.Exit(1)
}

// fileExists is a helper to check if a file exists.
func fileExists(f string) bool {
	_, err := os.Stat(f)
	return err == nil // this is wrong but whatever
}

// runServer the server, shutting down gracefully.
func runServer(ctx context.Context, server *http.Server) error {
	// ensure we wait for all tasks to finish.
	var (
		wg  sync.WaitGroup
		err error
	)

	// start the server in a goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		slog.InfoContext(ctx, "running HTTP server", "addr", server.Addr)
		err = server.ListenAndServe()
	}()

	// Wait for a shutdown signal in another goroutine,
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		slog.Info("shutting down gracefully")

		// 5 seconds should be plenty
		sd, ccl := context.WithTimeout(context.Background(), 5*time.Second)
		defer ccl()
		if err := server.Shutdown(sd); err != nil {
			slog.ErrorContext(ctx, "could not shut down gracefully", "error", err)
		}
	}()

	wg.Wait()

	return err
}

// ShortlinkNotFoundError is returned by Shortlinks.Long when a long link is not found.
type ShortlinkNotFoundError struct {
	Short string
}

// NewShortlinkNotFoundError returns a new ShortlinkNotFoundError.
func NewShortlinkNotFoundError(short string) *ShortlinkNotFoundError {
	return &ShortlinkNotFoundError{Short: short}
}

// Error implements the [error] interface.
func (e *ShortlinkNotFoundError) Error() string {
	return fmt.Sprintf("%s is unknown", e.Short)
}

// ShortlinkHandler is an HTTP handler that serves redirects for a [Shortlinks] implementation.
type ShortlinkHandler struct {
	Shortlinks
}

// ServeHTTP implements the [http.Handler] interface.
func (h ShortlinkHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	short := r.PathValue("short")
	if short == "health" {
		okHandler(w, r)
		return
	}

	long, err := h.Long(r.Context(), short)
	if err != nil {
		var nfe *ShortlinkNotFoundError
		if errors.As(err, &nfe) {
			http.NotFound(w, r)
			return
		}

		// Unknown error type.
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect the user.
	http.Redirect(w, r, long.String(), http.StatusFound)
}

// okHandler is a handler that always returns 200 OK
func okHandler(w http.ResponseWriter, r *http.Request) {
	h := w.Header()
	h.Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintln(w, http.StatusText(http.StatusOK))
}
