package mockfoo

import (
	"testing"

	proto "github.com/golang/protobuf/proto"
	assert "github.com/stretchr/testify/assert"
	pb "github.com/weathersource/go-gsrv/examples/foo/proto"
)

func TestInvokeSimple(t *testing.T) {
	s, err := NewServer()
	assert.NotNil(t, s)
	assert.Nil(t, err)

	s.AddData(
		"Bar",
		&pb.BarRequest{Baz: 1},
		&pb.BarResponse{Qux: "One"},
	)
	assert.Equal(t, 1, len(s.data))

	r, err := s.getData("Bar", &pb.BarRequest{Baz: 1})
	assert.NotNil(t, r)
	assert.Nil(t, err)

	s.Reset()
	assert.Equal(t, 0, len(s.data))
}

func TestInvokeNoRpc(t *testing.T) {
	s, err := NewServer()
	assert.NotNil(t, s)
	assert.Nil(t, err)

	r, err := s.getData("Foo", nil)
	assert.Nil(t, r)
	assert.NotNil(t, err)
}

func TestInvokeMissingData(t *testing.T) {
	s, err := NewServer()
	assert.NotNil(t, s)
	assert.Nil(t, err)

	s.AddData(
		"Bar",
		&pb.BarRequest{Baz: 1},
		&pb.BarResponse{Qux: "One"},
	)
	assert.Equal(t, 1, len(s.data))

	r, err := s.getData("Bar", &pb.BarRequest{Baz: 2})
	assert.Nil(t, r)
	assert.NotNil(t, err)

	s.Reset()
	assert.Equal(t, 0, len(s.data))
}

func TestInvokeAdjust(t *testing.T) {
	s, err := NewServer()
	assert.NotNil(t, s)
	assert.Nil(t, err)

	s.AddDataAdjust(
		"Bar",
		&pb.BarRequest{Baz: 1},
		&pb.BarResponse{Qux: "One"},
		func(wantReq proto.Message, gotReq proto.Message) proto.Message {
			wantReqAdj := *(wantReq.(*pb.BarRequest))
			wantReqAdj.Baz++
			return &wantReqAdj
		},
	)
	assert.Equal(t, 1, len(s.data))

	r, err := s.getData("Bar", &pb.BarRequest{Baz: 0})
	assert.NotNil(t, r)
	assert.Nil(t, err)

	s.Reset()
	assert.Equal(t, 0, len(s.data))
}
