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

type Queue struct {
	buf  []string
	size int
}

func InitQueue(size int) *Queue {
	buffer := make([]string, 0, size)
	return &Queue{
		buf:  buffer,
		size: 0,
	}
}

func (q *Queue) Push(elem string) {
	if q.size >= cap(q.buf) {
		q.buf = q.buf[1:]
		q.buf = append(q.buf, elem)
		return
	}
	q.buf = append(q.buf, elem)
	q.size++
}

// вспомогательная структура, содержит индекс строки , строку
// и поле isMatch - совпала ли с паттерном или нет
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
	//если требуется точное совпадение
	if f.fixed {
		for i, _ := range *s {
			//если нет требования игнорировать регистр
			if !f.ignorecase {
				//если необходимо исключать, возвращаем отрицание
				if !f.invert {
					//проверяем на наличие подстроки
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
			//проверяем на совпадение по паттерну
			if !f.ignorecase {
				matched, err = regexp.Match(f.rExp, []byte((*s)[i].s))
			} else {
				matched, err = regexp.Match(f.rExp, []byte(strings.ToLower((*s)[i].s)))
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "error compiling regex\n")
				(*s)[i].isMatch = false
			}
			//если требуется искключение
			if f.invert {
				matched = !matched
			}
			//у каждой строки метка, подходит под параметр или нет
			(*s)[i].isMatch = matched
		}
	}
}

// печать уже отфильтрованных значений
func printFiltered(s *[]stringLine, f Flags) {
	if f.after == 0 && f.before == 0 && f.context == 0 {
		//если не требуется буфер(флаги a, b, C)
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
		//если был флаг a,b или C
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
	//добавляем в выходной слайс все строки, которые имеют метку
	out := make([]string, 0, cap(*s))
	for _, v := range *s {
		if v.isMatch {
			//если было условие пронумеровать
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
	//если был флаг C, то a=b
	if f.context > 0 {
		cBefore = f.context
		cAfter = f.context
	} else {
		cBefore = f.before
		cAfter = f.after
	}
	//буфер строк до/после
	cBeforeBuf := InitQueue(cBefore)
	cAfterBuf := make([]string, 0, cAfter)
	//выходной
	out := make([]string, 0, cap(*s))
	//befIndex := 0
	afIndex := 0
	needAfter := false
	for _, v := range *s {
		//если нашли value и требуется добавить строки after
		if needAfter {
			//добавляем в буфер after
			if f.linenum {
				cAfterBuf = append(cAfterBuf, strconv.Itoa(v.index)+": "+v.s)
			} else {
				cAfterBuf = append(cAfterBuf, v.s)
			}
			afIndex++
			//когда добавили все строки по параметру after, добавляем их в выходной массив, чистим буфер, обнуляем индекс
			if afIndex == cAfter {
				afIndex = 0
				out = append(out, cAfterBuf...)
				cAfterBuf = make([]string, 0, cAfter)
				needAfter = v.isMatch
				if !needAfter {
					continue
				}
			}
		}
		//заполняем буфер before
		if !v.isMatch && !needAfter {
			//используем очередь
			if f.linenum {
				cBeforeBuf.Push(strconv.Itoa(v.index) + ": " + v.s)
			} else {
				cBeforeBuf.Push(v.s)
			}
		}
		//если нашли совпадение и не было требования печатать after
		if v.isMatch && !needAfter {
			//если был буфер before, добавляем в выходной слайс
			if len(cBeforeBuf.buf) > 0 {
				out = append(out, cBeforeBuf.buf...)
			}
			//добавляем строку которая совпала
			if f.linenum {
				out = append(out, strconv.Itoa(v.index)+": "+v.s)
			} else {
				out = append(out, v.s)
			}
			//обнуляем буфер before
			cBeforeBuf = InitQueue(cBefore)
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
	//читаем строки из файла в структуру
	strs := readFileIntoStruct(f)
	//отмечаем, попадает ли строка под заданный паттерн или нет
	markMatch(&strs, f)
	printFiltered(&strs, f)
}

func main() {
	//пример запуска
	//go run . -A 2 -B 2 -r 'lang.*' -n
	flags := flagsInit()
	StartSearching(*flags)
}
