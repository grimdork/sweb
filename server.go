package sweb

import (
	"context"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/Urethramancer/signor/log"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// Server structure.
type Server struct {
	sync.RWMutex
	sync.WaitGroup
	log.LogShortcuts
	http.Server

	web *chi.Mux

	staticpath string
}

const (
	// WEBHOST default
	WEBHOST = "127.0.0.1"
	// WEBPORT default
	WEBPORT = "15000"
	// WEBSTATIC default
	WEBSTATIC = "static"
)

// New server init. Reads settings from environment:
// WEBHOST - default 127.0.0.1
// WEBPORT - default 15000
// WEBSTATIC - default "./static/"
func New() *Server {
	srv := &Server{}
	srv.web = chi.NewRouter()

	// Logging
	srv.Logger = log.Default
	srv.L = log.Default.TMsg
	srv.E = log.Default.TErr

	// Routes
	srv.web.Use(
		middleware.RealIP,
		middleware.RequestID,
		srv.addLogger,
		AddHTMLHeaders,
	)

	srv.WebGet("/", srv.Static)
	srv.WebGets("/{page}", func(r chi.Router) {
		r.Get("/*", srv.Static)
		r.Options("/", Preflight)
	})

	return srv
}

// Start serving, reconfiguring from any changed environment variables.
func (srv *Server) Start() error {
	srv.Lock()
	defer srv.Unlock()

	srv.staticpath = getenv("WEBSTATIC", WEBSTATIC)
	addr := net.JoinHostPort(
		getenv("WEBHOST", WEBHOST),
		getenv("WEBPORT", WEBPORT),
	)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	srv.Add(1)
	srv.L("Starting web server on http://%s", addr)
	go func() {
		srv.Handler = srv.web
		err = srv.Serve(listener)

		if err != nil && err != http.ErrServerClosed {
			srv.E("Error running server: %s", err.Error())
			os.Exit(2)
		}
		srv.L("Stopped web server.")
		srv.Done()
	}()

	return nil
}

// Stop serving.
func (srv *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err := srv.Shutdown(ctx)
	if err != nil {
		srv.E("Shutdown error: %s", err.Error())
		os.Exit(2)
	}

	srv.Wait()
}

// func (ws *Server) wout(w http.ResponseWriter, s string) {
//         n, err := w.Write([]byte(s))
//         if err != nil {
//                 ws.E("Error: wrote %d bytes: %s", n, err.Error())
//         }
// }
