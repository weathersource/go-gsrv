package mockfoo

import (
	"testing"

	proto "google.golang.org/protobuf/proto"
	assert "github.com/stretchr/testify/assert"
	pb "github.com/weathersource/go-gsrv/examples/foo/proto"
)

func TestReset(t *testing.T) {
	s, err := NewServer()
	assert.NotNil(t, s)
	assert.Nil(t, err)

	s.AddRPC(
		&pb.BarRequest{Baz: 1},
		&pb.BarResponse{Qux: "One"},
	)
	assert.Equal(t, 1, len(s.resps))

	s.Reset()
	assert.Equal(t, 0, len(s.resps))
}

func TestPopRPCSimple(t *testing.T) {
	s, err := NewServer()
	assert.NotNil(t, s)
	assert.Nil(t, err)

	s.AddRPC(
		&pb.BarRequest{Baz: 1},
		&pb.BarResponse{Qux: "One"},
	)

	r, err := s.popRPC(&pb.BarRequest{Baz: 1})
	assert.NotNil(t, r)
	assert.Nil(t, err)
}

func TestPopRPCNoRpc(t *testing.T) {
	s, err := NewServer()
	assert.NotNil(t, s)
	assert.Nil(t, err)

	panicTest := func() {
		s.popRPC(nil)
	}
	assert.Panics(t, panicTest)
}

func TestPopRPCMissingData(t *testing.T) {
	s, err := NewServer()
	assert.NotNil(t, s)
	assert.Nil(t, err)

	s.AddRPC(
		&pb.BarRequest{Baz: 1},
		&pb.BarResponse{Qux: "One"},
	)

	r, err := s.popRPC(&pb.BarRequest{Baz: 2})
	assert.Nil(t, r)
	assert.NotNil(t, err)
}

func TestPopRPCAdjust(t *testing.T) {
	s, err := NewServer()
	assert.NotNil(t, s)
	assert.Nil(t, err)

	wantReq := &pb.BarRequest{Baz: 1}
	s.AddRPCAdjust(
		wantReq,
		&pb.BarResponse{Qux: "One"},
		func(gotReq proto.Message) {
			wantReq.Baz = gotReq.(*pb.BarRequest).Baz
		},
	)

	r, err := s.popRPC(&pb.BarRequest{Baz: 9})
	assert.NotNil(t, r)
	assert.Nil(t, err)
}
