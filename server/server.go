package server

type Server struct {
	filename string
}

func New(
	fn string,
) *Server {
	return &Server{
		filename: fn,
	}
}

func (s *Server) Run(stopChan chan error) {
	// do stuff
}
