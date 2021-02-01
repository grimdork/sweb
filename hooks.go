package sweb

// Hook allows for custom setup/teardown callbacks.
type Hook func() error

// AddStartHook takes a hook function to run before the server starts.
// You may also return an error, which makes the server fail launch.
// Hooks are executed in the order they were added.
func (srv *Server) AddStartHook(h Hook) {
	srv.starthooks = append(srv.starthooks, h)
}

// AddStopHook takes a hook function to run before the server stops.
// Hooks are executed in the order they were added.
func (srv *Server) AddStopHook(h Hook) {
	srv.stophooks = append(srv.stophooks, h)
}
