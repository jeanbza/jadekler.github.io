package server_test

import (
	"context"
	"errors"
	"testing"

	server "github.com/jadekler.github.io/example"
	grpc "google.golang.org/grpc"
)

type stubGreeter struct {
	sayHelloFn func(context.Context, *server.HelloRequest, ...grpc.CallOption) (*server.HelloReply, error)
}

func (s *stubGreeter) SayHello(ctx context.Context, req *server.HelloRequest, opts ...grpc.CallOption) (*server.HelloReply, error) {
	if s.sayHelloFn != nil {
		return s.sayHelloFn(ctx, req, opts...)
	}
	return nil, errors.New("SayHello is not set up with a stub")
}

func TestDoSomeWork(t *testing.T) {
	stub := &stubGreeter{}
	var sayHelloRequests []*server.HelloRequest
	stub.sayHelloFn = func(ctx context.Context, req *server.HelloRequest, opts ...grpc.CallOption) (*server.HelloReply, error) {
		sayHelloRequests = append(sayHelloRequests, req)
		return &server.HelloReply{}, nil
	}

	s := server.New(stub)
	if err := s.DoSomeWork(); err != nil {
		t.Fatal(err)
	}

	if got, want := 1, len(sayHelloRequests); got != want {
		t.Fatalf("got %d SayHello requests, want %d", got, want)
	}
	if got, want := sayHelloRequests[0].Name, "Jane Doe"; got != want {
		t.Errorf("got %s SayHelloRequest.Name, want %s", got, want)
	}
}
