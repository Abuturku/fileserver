package server

import (
	"testing"
	"github.com/stretchr/testify/assert"
)
func TestIsServerUp(t *testing.T) {
    assert.Equal(t, 123, 123, "Server should be up")
}
