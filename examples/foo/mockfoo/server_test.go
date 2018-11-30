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

	r, err := s.invoke("Bar", &pb.BarRequest{Baz: 1})
	assert.NotNil(t, r)
	assert.Nil(t, err)

	s.Reset()
	assert.Equal(t, 0, len(s.data))
}

func TestInvokeNoRpc(t *testing.T) {
	s, err := NewServer()
	assert.NotNil(t, s)
	assert.Nil(t, err)

	r, err := s.invoke("Foo", nil)
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

	r, err := s.invoke("Bar", &pb.BarRequest{Baz: 2})
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
		func(req proto.Message) proto.Message {
			thisReq := *(req.(*pb.BarRequest))
			thisReq.Baz++
			return &thisReq
		},
	)
	assert.Equal(t, 1, len(s.data))

	r, err := s.invoke("Bar", &pb.BarRequest{Baz: 0})
	assert.NotNil(t, r)
	assert.Nil(t, err)

	s.Reset()
	assert.Equal(t, 0, len(s.data))
}
