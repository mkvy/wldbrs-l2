package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOr(t *testing.T) {
	testI := make(chan interface{})
	close(testI)
	assert.Empty(t, or(testI))
}
