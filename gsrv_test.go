package gsrv

import (
	"testing"

	assert "github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	assert := assert.New(t)
	s, err := NewServer()
	if assert.Nil(err) {
		assert.Equal(Initialized, s.Status())
	}
}

func TestNewServerWithPort(t *testing.T) {
	assert := assert.New(t)
	s, err := NewServerWithPort(0)
	if assert.Nil(err) {
		assert.Equal(Initialized, s.Status())
	}
}

func TestNewServerWithPortError(t *testing.T) {
	assert := assert.New(t)
	s, err := NewServerWithPort(100000)
	assert.Nil(s)
	assert.NotNil(err)
}

func TestStart(t *testing.T) {
	assert := assert.New(t)
	s, err := NewServer()
	if assert.Nil(err) {
		s.Start()
		assert.Equal(Running, s.Status())
	}
}

func TestClose(t *testing.T) {
	assert := assert.New(t)
	s, err := NewServer()
	if assert.Nil(err) {
		s.Start()
		s.Close()
		assert.Equal(Closed, s.Status())
	}
}

func TestStatusString(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("Initialized", status(0).String())
	assert.Equal("9", status(9).String())
}
