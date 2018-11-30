package mockfoo

import (
	"testing"

	assert "github.com/stretchr/testify/assert"
	errors "github.com/weathersource/go-errors"
	pb "github.com/weathersource/go-gsrv/examples/foo/proto"
	context "golang.org/x/net/context"
)

func TestBarSuccess(t *testing.T) {
	s, err := NewServer()
	assert.NotNil(t, s)
	assert.Nil(t, err)

	s.AddData(
		"Bar",
		&pb.BarRequest{Baz: 1},
		&pb.BarResponse{Qux: "One"},
	)
	assert.Equal(t, 1, len(s.data))

	r, err := s.Bar(context.Background(), &pb.BarRequest{Baz: 1})
	assert.NotNil(t, r)
	assert.Nil(t, err)
}

func TestBarError(t *testing.T) {
	s, err := NewServer()
	assert.NotNil(t, s)
	assert.Nil(t, err)

	s.AddData(
		"Bar",
		&pb.BarRequest{Baz: 1},
		errors.NewUnknownError("An unknown error occurred."),
	)
	assert.Equal(t, 1, len(s.data))

	r, err := s.Bar(context.Background(), &pb.BarRequest{Baz: 1})
	assert.Nil(t, r)
	assert.NotNil(t, err)
}
