package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCutStrings(t *testing.T) {
	strs := []string{`aa bb cc`, `ee kk qq`, `ee kk`}
	flags := &Flags{delimiter: " ", fields: []int{1, 3}}
	outExpected := []string{`aa cc`, `ee qq`, `ee`}
	require.Equal(t, cutStrings(strs, *flags), outExpected)
}

func TestCutStringsSep(t *testing.T) {
	strs := []string{`aa bb cc`, `ee kk qq`, `ee`}
	flags := &Flags{delimiter: " ", fields: []int{1, 3}, sep: true}
	outExpected := []string{`aa cc`, `ee qq`}
	require.Equal(t, cutStrings(strs, *flags), outExpected)
}
