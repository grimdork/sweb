package sweb

import (
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi"
)

// WebGet adds a GET route matching the specified pattern.
func (srv *Server) WebGet(pattern string, handler http.HandlerFunc) {
	srv.web.Get(pattern, handler)
	srv.L("Added GET route for %s", pattern)
}

// WebGets adds one or more GET routes to the specified pattern.
func (srv *Server) WebGets(pattern string, fn func(r chi.Router)) {
	srv.web.Route(pattern, fn)
	srv.L("Added GET routes for %s", pattern)
}

// Route adds more sub-routes to a chi route.
func (srv *Server) Route(pattern string, fn func(r chi.Router)) chi.Router {
	return srv.web.Route(pattern, fn)
}

// Static page serving.
func (srv *Server) Static(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Path
	if page == "/" {
		page = "/index.html"
	}

	srv.ServeFile(w, r, page)
}

// ServeFile serves a file from the WEBSTATIC path.
// NOTE: The server needs to be reloaded if the environment somehow changes.
func (srv *Server) ServeFile(w http.ResponseWriter, r *http.Request, name string) {
	fn := filepath.Join(srv.staticpath, name)
	srv.L("Serving file '%s' from filename '%s'", name, fn)
	f, err := os.Open(fn)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	defer f.Close()
	ext := filepath.Ext(fn)
	if ext != "" {
		w.Header().Set("Content-Type", mime.TypeByExtension(ext))
	} else {
		w.Header().Set("Content-Type", mime.TypeByExtension(".txt"))
	}

	http.ServeContent(w, r, name, time.Time{}, f)
}
