package sweb

// StartHook allows for custom setup and cancelling the setup process.
type StartHook func() error

// StopHook allows for custom teardown.
type StopHook func()

// AddStartHook takes a hook function to run before the server starts.
// You may also return an error, which makes the server fail launch.
// Hooks are executed in the order they were added.
func (srv *Server) AddStartHook(h StartHook) {
	srv.starthooks = append(srv.starthooks, h)
}

// AddStopHook takes a hook function to run before the server stops.
// Hooks are executed in the order they were added.
func (srv *Server) AddStopHook(h StopHook) {
	srv.stophooks = append(srv.stophooks, h)
}
