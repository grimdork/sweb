# sweb
Serve webstuff.

## What?
This is a quickstart package to get a web server app running, waiting on a port and IP address specified in the environment, and serving static files from a similarly specified directory.

This is geared at things running in typical Docker, AWS etc. setups behind a reverse proxy which handles the domain and certificates. If you need S3, GCP or other special backend storage support, write a route for it. See chi documentation for more on how routing works.

## How do I try it?
The minimal example:

```go
srv := sweb.New()
srv.Start()
// do stuff, wait for CTRL-C, whatever
srv.Stop()
```

This will launch a web server bound to the IP address 127.0.0.1, port 15000, loading static files from the folder `static` in the working path of the program.

To configure these settings, use the following environment variables:
```
WEBHOST
WEBPORT
WEBSTATIC
```

# How do I customise the root path?
Embed the Server struct into a struct of your own, run srv.Init(), optionally srv.InitMiddleware() and add routes manually from there with WebGet(), WebGets() and Route().

Example:
```go
type MyServer struct {
	sweb.Server
}
…
srv:=&MyServer{}
// Add router and logger
srv.Init()
// Add four pieces of middleware suitable for HTTP
srv.InitMiddleware()
srv.WebGet("/", srv.Static)
srv.WebGets("/{page}", func(r chi.Router) {
	r.Get("/*", srv.Static)
	r.Options("/", sweb.Preflight)
})
…
srv.Stop()
```

This works like the simpler example, setting things up like the defaults.

Example setup for a `/api` route alongside the default handler:
```go
srv.Route("/api", func(r chi.Router) {
	// API middleware
	r.Use(
		middleware.NoCache, // chi
		middleware.RealIP, // chi
		sweb.AddCORS,
		middleware.Timeout(time.Second*10), // chi
	)
	// TODO: Insert r.NotFound(w) call here if desirable
	r.Options("/", sweb.Preflight)
	// Longer endpoints first, otherwise the "/" route will be used for everything
	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello " + r.RemoteAddr))
	})
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("v1"))
	})
})
```

## Middleware
These middlware functions are available for use, in addition to everything found in chi:

- AddJSONHeaders: For typical REST endpoints
- AddHTMLHeaders: For regular browser-friendly pages
- Preflight: Sets up options for REST calls
- AddCORS: Allows REST calls from other domains
- AddSecureHeaders: Forces secure pages when behind a HTTPS reverse proxy
