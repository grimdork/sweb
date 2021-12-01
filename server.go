package sweb

import (
	"context"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	ll "github.com/grimdork/loglines"
)

// Server structure.
type Server struct {
	sync.RWMutex
	sync.WaitGroup
	http.Server

	web        *chi.Mux
	starthooks []StartHook
	stophooks  []StopHook

	staticpath string
}

const (
	// WEBHOST default
	WEBHOST = "0.0.0.0"
	// WEBPORT default
	WEBPORT = "80"
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

// Use passes middleware to the root web route.
func (srv *Server) Use(middlewares ...func(http.Handler) http.Handler) {
	srv.web.Use(middlewares...)
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
	ll.Msg("Starting web server on http://%s", addr)
	go func() {
		srv.Handler = srv.web
		err = srv.Serve(listener)

		if err != nil && err != http.ErrServerClosed {
			ll.Err("Error running server: %s", err.Error())
			os.Exit(2)
		}
		ll.Msg("Stopped web server.")
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
		ll.Err("Shutdown error: %s", err.Error())
	}

	for _, cb := range srv.stophooks {
		cb()
	}
	srv.Wait()
}
