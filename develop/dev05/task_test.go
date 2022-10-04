package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResultGrepFilter(t *testing.T) {
	flags := &Flags{rExp: `orl.*`}
	testSlice := []string{"world"}
	r := &resultGrepFilter{}
	s := make([]stringLine, 0)
	s = append(s, stringLine{s: "golang", index: 0}, stringLine{s: "world", index: 1, isMatch: false}, stringLine{s: "foodsea", index: 2})
	markMatch(&s, *flags)
	assert.Equal(t, r.resFilter(&s, *flags), testSlice)
}

func TestResultGrepFilterBuf(t *testing.T) {
	flags := &Flags{rExp: `orl.*`, after: 1, before: 1, linenum: true}
	testSlice := []string{"0: golang", "1: world", "2: foodsea"}
	r := &resultGrepFilterBuf{}
	s := make([]stringLine, 0)
	s = append(s, stringLine{s: "golang", index: 0}, stringLine{s: "world", index: 1, isMatch: false}, stringLine{s: "foodsea", index: 2})
	markMatch(&s, *flags)
	assert.Equal(t, r.resFilter(&s, *flags), testSlice)
}

func TestResultGrepCount(t *testing.T) {
	flags := &Flags{rExp: `lang.*`, count: true, linenum: true}
	r := &resultGrepFilter{}
	s := make([]stringLine, 0)
	s = append(s, stringLine{s: "golang", index: 0}, stringLine{s: "world", index: 1, isMatch: false}, stringLine{s: "language", index: 2})
	markMatch(&s, *flags)
	assert.Equal(t, len(r.resFilter(&s, *flags)), 2)
}
