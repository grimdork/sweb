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

	web        *chi.Mux
	starthooks []Hook
	stophooks  []Hook

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
	srv.Init()
	srv.InitMiddleware()
	srv.InitRouter()
	return srv
}

// Init sets up the basics, but no routes.
func (srv *Server) Init() {
	srv.web = chi.NewRouter()

	// Logging
	srv.Logger = log.Default
	srv.L = log.Default.TMsg
	srv.E = log.Default.TErr

	// Default timeouts
	srv.Server.IdleTimeout = time.Second * 30
	srv.Server.ReadTimeout = time.Second * 30
	srv.Server.ReadHeaderTimeout = time.Second * 5
	srv.Server.WriteTimeout = time.Second * 30
}

// InitMiddleware sets up basic middleware on the root web route.
// These are RealIP and RequestID from chi, a logger for visits and HTML headers.
func (srv *Server) InitMiddleware() {
	srv.web.Use(
		middleware.RealIP,
		middleware.RequestID,
		srv.addLogger,
		AddHTMLHeaders,
	)
}

// InitRouter creates the default root router which loads files from the WEBSTATIC path.
func (srv *Server) InitRouter() {
	srv.WebGet("/", srv.Static)
	srv.WebGets("/{page}", func(r chi.Router) {
		r.Get("/*", srv.Static)
		r.Options("/", Preflight)
	})
}

// Start serving, reconfiguring from any changed environment variables.
func (srv *Server) Start() error {
	for _, cb := range srv.starthooks {
		err := cb()
		if err != nil {
			return err
		}
	}

	srv.Lock()
	defer srv.Unlock()

	srv.staticpath = Getenv("WEBSTATIC", WEBSTATIC)
	addr := net.JoinHostPort(
		Getenv("WEBHOST", WEBHOST),
		Getenv("WEBPORT", WEBPORT),
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
	}

	for _, cb := range srv.stophooks {
		cb()
	}
	srv.Wait()
}
