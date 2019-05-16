package gsrv

import (
	"fmt"
	"log"
	"net"
	"regexp"
	"strconv"
	"sync"

	errors "github.com/weathersource/go-errors"
	grpc "google.golang.org/grpc"
)

type status uint8

const (
	// Initialized Server status means the server has been created, but not started.
	Initialized status = iota
	// Running Server status means the server has been started.
	Running
	// Closed Server status means the server has been closed
	Closed
	// ClosedWithError Server status means the server was closed because of an error.
	ClosedWithError
)

func (s status) String() string {
	name := []string{"Initialized", "Running", "Closed", "ClosedWithError"}
	i := uint8(s)
	switch {
	case i <= uint8(Closed):
		return name[i]
	default:
		return strconv.Itoa(int(uint8(i)))
	}
}

// Server is an in-process gRPC server, listening on a system-chosen port on
// the local loopback interface. Servers are for testing only and are not
// intended to be used in production code.
//
// To create a server, make a new Server, register your handlers, then call
// Start:
//
//  srv, err := NewServer()
//  ...
//  mypb.RegisterMyServiceServer(srv.Gsrv, &myHandler)
//  ....
//  srv.Start()
//
// Clients should connect to the server with no security:
//
//  conn, err := grpc.Dial(srv.Addr, grpc.WithInsecure())
//  ...
type Server struct {
	Addr   string
	Port   int
	l      net.Listener
	Gsrv   *grpc.Server
	status status
}

// NewServer creates a new Server. The server will be listening for gRPC connections
// at the address named by the Addr field, without TLS.
func NewServer(opts ...grpc.ServerOption) (*Server, error) {
	return NewServerWithPort(0, opts...)
}

// NewServerWithPort creates a new Server at a specific port. The server will be listening
// for gRPC connections at the address named by the Addr field, without TLS.
func NewServerWithPort(port int, opts ...grpc.ServerOption) (*Server, error) {
	if port > 65535 {
		return nil, errors.NewInvalidArgumentError("port must be in [0,65535].")
	}
	l, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return nil, errors.NewUnknownError("Failed to listen.", err)
	}
	s := &Server{
		Addr: l.Addr().String(),
		Port: parsePort(l.Addr().String()),
		l:    l,
		Gsrv: grpc.NewServer(opts...),
	}
	return s, nil
}

// Start causes the server to start accepting incoming connections.
// Call Start after registering handlers.
func (s *Server) Start() {
	if s.status == Initialized {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			s.status = Running
			wg.Done()
			err := s.Gsrv.Serve(s.l)
			if err != nil {
				s.status = ClosedWithError
				log.Println(errors.NewUnknownError("Failed to serve.", err))

			}
		}()
		wg.Wait()
	}
}

// Close shuts down the server.
func (s *Server) Close() {
	if s.status == Running {
		s.Gsrv.Stop()
		s.l.Close()
	}
	if s.status < Closed {
		s.status = Closed
	}
}

// Status returns the server status: Initialized, Running, Closed, or
// ClosedWithError
func (s *Server) Status() status {
	return s.status
}

var portParser = regexp.MustCompile(`:[0-9]+`)

func parsePort(addr string) int {
	res := portParser.FindAllString(addr, -1)
	stringPort := res[0][1:] // strip the :
	p, _ := strconv.ParseInt(stringPort, 10, 32)
	return int(p)
}
