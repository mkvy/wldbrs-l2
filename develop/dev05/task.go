package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

/*
=== Утилита grep ===

Реализовать утилиту фильтрации (man grep)

Поддержать флаги:
-A - "after" печатать +N строк после совпадения
-B - "before" печатать +N строк до совпадения
-C - "context" (A+B) печатать ±N строк вокруг совпадения
-c - "count" (количество строк)
-i - "ignore-case" (игнорировать регистр)
-v - "invert" (вместо совпадения, исключать)
-F - "fixed", точное совпадение со строкой, не паттерн
-n - "line num", печатать номер строки

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type stringLine struct {
	s       string
	isMatch bool
	index   int
}

type Flags struct {
	after      int
	before     int
	context    int
	count      bool
	ignorecase bool
	invert     bool
	fixed      bool
	linenum    bool
	rExp       string
	filename   string
}

func markMatch(s *[]stringLine, f Flags) {
	if f.fixed {
		for i, _ := range *s {
			if !f.ignorecase {
				if !f.invert {
					(*s)[i].isMatch = strings.Contains((*s)[i].s, f.rExp)
				} else {
					(*s)[i].isMatch = !strings.Contains((*s)[i].s, f.rExp)
				}
			} else {
				if !f.invert {
					(*s)[i].isMatch = strings.Contains(strings.ToLower((*s)[i].s), f.rExp)
				} else {
					(*s)[i].isMatch = !strings.Contains(strings.ToLower((*s)[i].s), f.rExp)
				}
			}
		}
	} else {
		for i, _ := range *s {
			var matched bool
			var err error
			if !f.ignorecase {
				matched, err = regexp.Match(f.rExp, []byte((*s)[i].s))
			} else {
				matched, err = regexp.Match(f.rExp, []byte(strings.ToLower((*s)[i].s)))
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "error compiling regex\n")
				(*s)[i].isMatch = false
			}
			if f.invert {
				matched = !matched
			}
			(*s)[i].isMatch = matched
		}
	}
}

func printFiltered(s *[]stringLine, f Flags) {
	if f.after == 0 && f.before == 0 && f.context == 0 {
		r := &resultGrepFilter{}
		res := r.resFilter(s, f)
		if !f.count {
			for _, v := range res {
				fmt.Println(v)
			}
		} else {
			fmt.Println(len(res))
		}
	} else {
		r := &resultGrepFilterBuf{}
		res := r.resFilter(s, f)
		for _, v := range res {
			fmt.Println(v)
		}
	}
}

type GetGrepRes interface {
	resFilter(s *[]stringLine, f Flags) []string
}

type resultGrepFilter struct{}

func (r *resultGrepFilter) resFilter(s *[]stringLine, f Flags) []string {
	out := make([]string, 0, cap(*s))
	for _, v := range *s {
		if v.isMatch {
			if f.linenum {
				out = append(out, strconv.Itoa(v.index)+": "+v.s)
			} else {
				out = append(out, v.s)
			}
		}
	}
	return out
}

type resultGrepFilterBuf struct{}

func (r *resultGrepFilterBuf) resFilter(s *[]stringLine, f Flags) []string {
	var cBefore, cAfter int
	if f.context > 0 {
		cBefore = f.context
		cAfter = f.context
	} else {
		cBefore = f.before
		cAfter = f.after
	}
	cBeforeBuf := make([]string, 0, cBefore)
	cAfterBuf := make([]string, 0, cAfter)
	out := make([]string, 0, cap(*s))
	befIndex := 0
	afIndex := 0
	needAfter := false
	for _, v := range *s {
		if needAfter {
			if f.linenum {
				cAfterBuf = append(cAfterBuf, strconv.Itoa(v.index)+": "+v.s)
			} else {
				cAfterBuf = append(cAfterBuf, v.s)
			}
			afIndex++
			if afIndex == cAfter {
				afIndex = 0
				out = append(out, cAfterBuf...)
				cAfterBuf = make([]string, 0, cAfter)
				needAfter = false
			}
		}
		if !v.isMatch && !needAfter && befIndex < cBefore {
			if f.linenum {
				cBeforeBuf = append(cBeforeBuf, strconv.Itoa(v.index)+": "+v.s)
			} else {
				cBeforeBuf = append(cBeforeBuf, v.s)
			}
			befIndex++
			if befIndex == cBefore {
				befIndex = 0
			}
		}
		if v.isMatch && !needAfter {
			if len(cBeforeBuf) > 0 {
				out = append(out, cBeforeBuf...)
			}
			if f.linenum {
				out = append(out, strconv.Itoa(v.index)+": "+v.s)
			} else {
				out = append(out, v.s)
			}
			befIndex = 0
			cBeforeBuf = make([]string, 0, cBefore)
			if cAfter > 0 {
				needAfter = true
			}
		}
	}
	return out
}

func flagsInit() *Flags {
	f := new(Flags)
	flag.IntVar(&f.after, "A", 0, "n strings after match")
	flag.IntVar(&f.before, "B", 0, "n strings before match")
	flag.IntVar(&f.context, "C", 0, "n strings after/before match")
	flag.BoolVar(&f.count, "c", false, "strings needs to be count")
	flag.BoolVar(&f.ignorecase, "i", false, "ignore case")
	flag.BoolVar(&f.invert, "v", false, "invert")
	flag.BoolVar(&f.fixed, "F", false, "not regex string")
	flag.BoolVar(&f.linenum, "n", false, "print line num")
	flag.StringVar(&f.filename, "file", "./file.txt", "read file name")
	flag.StringVar(&f.rExp, "r", `.`, "regex")
	flag.Parse()
	return f
}

func readFileIntoStruct(f Flags) []stringLine {
	file, err := os.Open(f.filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	strs := make([]stringLine, 0)
	cntr := 1
	for scanner.Scan() {
		strs = append(strs, stringLine{s: scanner.Text(), index: cntr})
		cntr++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return strs
}

func StartSearching(f Flags) {
	strs := readFileIntoStruct(f)
	markMatch(&strs, f)
	printFiltered(&strs, f)
}

func main() {
	flags := flagsInit()
	StartSearching(*flags)
}
