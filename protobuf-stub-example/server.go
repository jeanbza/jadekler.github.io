package server

import context "context"

type Server struct {
	gc GreeterClient
}

func New(gc GreeterClient) *Server {
	return &Server{gc: gc}
}

func (s *Server) DoSomeWork() error {
	ctx := context.Background()
	resp, err := s.gc.SayHello(ctx, &HelloRequest{Name: "Jane Doe"})
	if err != nil {
		return err
	}
	_ = resp // TODO: Do something with resp.
	return nil
}
