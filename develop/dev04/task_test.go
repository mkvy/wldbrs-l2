package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetAnagramSet(t *testing.T) {
	strs := []string{"пятка", "тяпкА", "слиток", "столик", "СТОЛИК", "Листок", "ТЯпка", "пятак"}
	mapCheck := map[string][]string{"пятка": {"пятак", "пятка", "тяпка"}, "слиток": {"листок", "слиток", "столик"}}
	require.Equal(t, GetAnagramSet(strs), mapCheck)
}

func TestUniqueStrs(t *testing.T) {
	strs := []string{"a", "b", "c", "b", "c", "b"}
	req := []string{"a", "b", "c"}
	require.Equal(t, uniqueStrs(strs), req)
}

func TestAnagram(t *testing.T) {
	str1 := "пятка"
	str2 := "тяпка"
	require.True(t, isAnagram(str1, str2))
	str2 = "тяпкат"
	require.False(t, isAnagram(str1, str2))
	str1, str2 = str2, str1
	require.False(t, isAnagram(str1, str2))
}
