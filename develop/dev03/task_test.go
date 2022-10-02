package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestQuickSortDefault(t *testing.T) {
	flags := &flagsCmd{}
	assert.Equal(t, quickSort([]string{"ba bc ce", "ab jk kk", "dd"}, *flags), []string{"ab jk kk", "ba bc ce", "dd"})
}

func TestQuickSortByColumn(t *testing.T) {
	flags := &flagsCmd{k: 2}
	assert.Equal(t, quickSort([]string{"ba bc ce", "ab jk kk", "dd er", "aa cc kk"}, *flags), []string{"ba bc ce", "aa cc kk", "dd er", "ab jk kk"})
}

func TestQuickSortByNumber(t *testing.T) {
	flags := &flagsCmd{n: true}
	assert.Equal(t, quickSort([]string{"4 5 6", "1 2 3", "2 2 3"}, *flags), []string{"1 2 3", "2 2 3", "4 5 6"})
}

func TestQuickSortByNumberColumn(t *testing.T) {
	flags := &flagsCmd{n: true, k: 3}
	assert.Equal(t, quickSort([]string{"4 5 9", "1 2 5", "2 2 0"}, *flags), []string{"2 2 0", "1 2 5", "4 5 9"})
}

func TestQuickSortBySuffixNum(t *testing.T) {
	flags := &flagsCmd{h: true}
	assert.Equal(t, quickSort([]string{"1066.5B", "10.1G", "1.0K"}, *flags), []string{"1.0K", "1066.5B", "10.1G"})
}

func TestQuickSortByMonth(t *testing.T) {
	flags := &flagsCmd{M: true}
	assert.Equal(t, quickSort([]string{"jul", "aug", "jan"}, *flags), []string{"jan", "jul", "aug"})
}

func TestQuickSortByMonthReverse(t *testing.T) {
	flags := &flagsCmd{M: true, r: true}
	assert.Equal(t, quickSort([]string{"jul", "aug", "jan"}, *flags), []string{"aug", "jul", "jan"})
}

func TestQuickSortUnique(t *testing.T) {
	flags := &flagsCmd{u: true}
	assert.Equal(t, quickSort([]string{"ee", "aa", "ee"}, *flags), []string{"aa", "ee"})
}
func TestQuickSortUniqueReverse(t *testing.T) {
	flags := &flagsCmd{u: true, r: true}
	assert.Equal(t, quickSort([]string{"ee", "aa", "ee"}, *flags), []string{"ee", "aa"})
}

func TestQuickSortCheck(t *testing.T) {
	flags := &flagsCmd{}
	assert.True(t, checkSort([]string{"aa", "ab", "bc"}, *flags), true)
}
