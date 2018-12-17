// Package mockfoo mocks a server for the Foo service.
package mockfoo

import (
	"fmt"

	proto "github.com/golang/protobuf/proto"
	errors "github.com/weathersource/go-errors"
	gsrv "github.com/weathersource/go-gsrv"
	pb "github.com/weathersource/go-gsrv/examples/foo/proto"
)

// MockServer satisfies the pb.FooServer interface
type MockServer struct {
	pb.FooServer // <<< VALIDATE CURRENT SERVICE SERVER

	Addr string

	reqItems []reqItem
	resps    []interface{}
}

type reqItem struct {
	wantReq proto.Message
	adjust  func(gotReq proto.Message)
}

// NewServer configures and returns a MockServer
func NewServer() (*MockServer, error) {
	srv, err := gsrv.NewServer()
	if err != nil {
		return nil, err
	}
	mock := &MockServer{Addr: srv.Addr}
	pb.RegisterFooServer(srv.Gsrv, mock) // <<< VALIDATE CURRENT SERVICE REGISTRY
	srv.Start()
	return mock, nil
}

// Reset returns the MockServer to an empty state.
func (s *MockServer) Reset() {
	s.reqItems = nil
	s.resps = nil
}

// AddRPC adds a (request, response) pair to the server's list of expected
// interactions. The server will compare the incoming request with wantReq
// using proto.Equal. The response can be a message or an error.
//
// For the Listen RPC, resp should be a []interface{}, where each element
// is either ListenResponse or an error.
//
// Passing nil for wantReq disables the request check.
func (s *MockServer) AddRPC(wantReq proto.Message, resp interface{}) {
	s.AddRPCAdjust(wantReq, resp, nil)
}

// AddRPCAdjust is like AddRPC, but accepts a function that can be used
// to tweak the requests before comparison, for example to adjust for
// randomness.
func (s *MockServer) AddRPCAdjust(wantReq proto.Message, resp interface{}, adjust func(gotReq proto.Message)) {
	s.reqItems = append(s.reqItems, reqItem{wantReq, adjust})
	s.resps = append(s.resps, resp)
}

// popRPC compares the request with the next expected (request, response) pair.
// It returns the response, or an error if the request doesn't match what
// was expected or there are no expected rpcs.
func (s *MockServer) popRPC(gotReq proto.Message) (interface{}, error) {
	if len(s.reqItems) == 0 {
		panic("mocks.popRPC: Out of RPCs.")
	}
	ri := s.reqItems[0]
	s.reqItems = s.reqItems[1:]
	if ri.wantReq != nil {
		if ri.adjust != nil {
			ri.adjust(gotReq)
		}

		if !proto.Equal(gotReq, ri.wantReq) {
			return nil, errors.NewInternalError(fmt.Sprintf("mocks.popRPC: Bad request\ngot:  %T\n%s\nwant: %T\n%s",
				gotReq, proto.MarshalTextString(gotReq),
				ri.wantReq, proto.MarshalTextString(ri.wantReq)))
		}
	}
	resp := s.resps[0]
	s.resps = s.resps[1:]
	if err, ok := resp.(error); ok {
		return nil, err
	}
	return resp, nil
}
