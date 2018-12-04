// Package mockfoo mocks a server for the Foo service.
package mockfoo

import (
	proto "github.com/golang/protobuf/proto"
	errors "github.com/weathersource/go-errors"
	gsrv "github.com/weathersource/go-gsrv"
	pb "github.com/weathersource/go-gsrv/examples/foo/proto"
)

// MockServer satisfies the pb.FooServer interface
type MockServer struct {
	pb.FooServer // <<< VALIDATE CURRENT SERVICE SERVER
	Addr         string
	data         map[string][]datum
}

type datum struct {
	req    proto.Message
	adjust AdjustFunc
	resp   interface{}
}

// AdjustFunc is the signature for functions that adjust requests before querying matching data.
// Ensure inside AdjustFunc functions, you modify and return a copy of the original requests,
// not pointers to the modified originals. This allows for running AdjustFunc multiple times
// against multiple possible results
type AdjustFunc func(wantReq proto.Message, gotReq proto.Message) (wantReqAdj proto.Message)

// NewServer configures and returns a MockServer
func NewServer() (*MockServer, error) {
	srv, err := gsrv.NewServer()
	if err != nil {
		return nil, err
	}
	mock := &MockServer{Addr: srv.Addr}
	mock.data = make(map[string][]datum)
	pb.RegisterFooServer(srv.Gsrv, mock) // <<< VALIDATE CURRENT SERVICE REGISTRY
	srv.Start()
	return mock, nil
}

// Reset returns the MockServer to an empty state.
func (s *MockServer) Reset() {
	s.data = make(map[string][]datum)
}

// AddData adds a (request, response) pair for a given rpc function name to the
// server's list of expected interactions. The server will compare the incoming
// request with req using proto.Equal. The response can be a message or an
// error.
//
// For the Listen RPC, resp should be a interface{}, where the element
// is either ListenResponse or an error.
//
// Passing nil for req disables the request check.
func (s *MockServer) AddData(rpc string, req proto.Message, resp interface{}) {
	s.AddDataAdjust(rpc, req, resp, nil)
}

// AddDataAdjust is like AddData, but accepts a function that can be used
// to tweak the requests before comparison, for example to adjust for
// randomness.
func (s *MockServer) AddDataAdjust(rpc string, req proto.Message, resp interface{}, adjust AdjustFunc) {
	_, ok := s.data[rpc]
	if !ok {
		s.data[rpc] = []datum{}
	}
	s.data[rpc] = append(s.data[rpc], datum{req, adjust, resp})
}

// getData compares the request with the next expected (request, response) pair.
// It returns the response, or an error if the request doesn't match what
// was expected or there are no expected data.
func (s *MockServer) getData(rpc string, wantReq proto.Message) (interface{}, error) {
	_, ok := s.data[rpc]
	if !ok {
		return nil, errors.NewNotFoundError("The RPC could not be found.")
	}

	for _, got := range s.data[rpc] {

		if got.req != nil {

			// Handle request adjustments with proper scope
			wantReqAdj := wantReq
			if got.adjust != nil {
				wantReqAdj = got.adjust(wantReq, got.req)
			}

			if !proto.Equal(wantReqAdj, got.req) {
				continue
			}
		}
		if err, ok := got.resp.(error); ok {
			return nil, err
		}
		return got.resp, nil
	}

	return nil, errors.NewNotFoundError("The requested response object was not found.")
}
