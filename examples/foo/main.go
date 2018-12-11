package main

import (
	"fmt"
	"log"

	errors "github.com/weathersource/go-errors"
	mockfoo "github.com/weathersource/go-gsrv/examples/foo/mockfoo"
	pb "github.com/weathersource/go-gsrv/examples/foo/proto"
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

func main() {
	// set up mock server
	srv, err := mockfoo.NewServer()
	if err != nil {
		log.Fatalf("did not serve: %v", err)
	}

	// populate server
	srv.AddRPC(
		&pb.BarRequest{Baz: 1},
		&pb.BarResponse{Qux: "One"},
	)
	srv.AddRPC(
		&pb.BarRequest{Baz: 2},
		errors.NewNotFoundError("An unknown error occurred."),
	)

	// Set up a connection to the server.
	conn, err := grpc.Dial(srv.Addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewFooClient(conn)

	r, err := c.Bar(context.Background(), &pb.BarRequest{Baz: 1})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(r)
	}

	r2, err := c.Bar(context.Background(), &pb.BarRequest{Baz: 2})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(r2)
	}
}
