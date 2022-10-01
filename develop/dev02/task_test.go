package main

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestUnpackStr1(t *testing.T) {
	str := `a4bc2d5e`
	strOut, err := UnpackStr(str)
	assert.Equal(t, strOut, `aaaabccddddde`)
	assert.Equal(t, err, nil)
}

func TestUnpackStr2(t *testing.T) {
	str := `abcd`
	strOut, err := UnpackStr(str)
	assert.Equal(t, strOut, `abcd`)
	assert.Equal(t, err, nil)
}

func TestUnpackStr3(t *testing.T) {
	str := `45`
	strOut, err := UnpackStr(str)
	assert.Equal(t, strOut, ``)
	assert.EqualError(t, err, "wrong string")
}

func TestUnpackStr4(t *testing.T) {
	str := ``
	strOut, err := UnpackStr(str)
	assert.Equal(t, strOut, ``)
	assert.EqualError(t, err, "empty string")
}

func TestUnpackStr5(t *testing.T) {
	str := `qwe\4\5`
	strOut, err := UnpackStr(str)
	assert.Equal(t, strOut, `qwe45`)
	assert.Equal(t, err, nil)
}

func TestUnpackStr6(t *testing.T) {
	str := `qwe\45`
	strOut, err := UnpackStr(str)
	assert.Equal(t, strOut, `qwe44444`)
	assert.Equal(t, err, nil)
}

func TestUnpackStr7(t *testing.T) {
	str := `qwe\\5`
	strOut, err := UnpackStr(str)
	assert.Equal(t, strOut, `qwe\\\\\`)
	assert.Equal(t, err, nil)
}

func TestIsSymbol(t *testing.T) {
	str := "A1_"
	rs := []rune(str)
	assert.Equal(t, isSymbol(rs[0]), true)
	assert.Equal(t, isSymbol(rs[1]), false)
	assert.Equal(t, isSymbol(rs[2]), true)
}

func TestIsDigit(t *testing.T) {
	str := "0a9"
	rs := []rune(str)
	assert.Equal(t, isDigit(rs[0]), true)
	assert.Equal(t, isDigit(rs[1]), false)
	assert.Equal(t, isDigit(rs[2]), true)
}

func TestWriteMultipleRunes(t *testing.T) {
	sb := strings.Builder{}
	rs := []rune("a")
	sb = writeMultipleRunes(&sb, rs[0], 4)
	assert.Equal(t, sb.String(), `aaa`)
}
