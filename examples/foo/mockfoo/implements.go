package mockfoo

import (
	pb "github.com/weathersource/go-gsrv/examples/foo/proto"
	context "golang.org/x/net/context"
)

// Bar implements the FooServer Bar method
func (s *MockServer) Bar(_ context.Context, req *pb.BarRequest) (*pb.BarResponse, error) {
	res, err := s.popRPC(req)
	if err != nil {
		return nil, err
	}
	return res.(*pb.BarResponse), nil
}
