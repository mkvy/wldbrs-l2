package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

/*
=== Утилита sort ===

Отсортировать строки (man sort)
Основное

# Поддержать ключи

-k — указание колонки для сортировки
-n — сортировать по числовому значению
-r — сортировать в обратном порядке
-u — не выводить повторяющиеся строки

# Дополнительное

# Поддержать ключи

-M — сортировать по названию месяца
-b — игнорировать хвостовые пробелы
-c — проверять отсортированы ли данные
-h — сортировать по числовому значению с учётом суффиксов

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/
var (
	Months map[string]byte = map[string]byte{"jan": 1, "feb": 2, "mar": 3, "apr": 4, "may": 5, "jun": 6, "jul": 7, "aug": 8, "sep": 9, "oct": 10, "nov": 11, "dec": 12}
	//степени 1000
	Numsufx map[string]float64 = map[string]float64{"b": 0, "k": 1, "m": 2, "g": 3, "t": 4, "p": 5}
)

func quickSort(arr []string, flags flagsCmd) []string {
	res := make([]string, len(arr))
	copy(res, arr)
	if flags.c {
		isSorted := checkSort(arr, flags)
		if isSorted {
			fmt.Fprintf(os.Stdout, "Strings is sorted!\n")
			return res
		}
	}
	if flags.k <= 0 && !flags.n && !flags.M && !flags.h {
		quickSortAlg(res, 0, len(res)-1, flags, &CompareByDefault{})
	}
	if flags.k > 0 {
		quickSortAlg(res, 0, len(res)-1, flags, &CompareByColumn{})
	}
	if flags.n && flags.k <= 0 {
		quickSortAlg(res, 0, len(res)-1, flags, &CompareByInt{})
	}
	if flags.M {
		quickSortAlg(res, 0, len(res)-1, flags, &CompareByMonth{})
	}
	if flags.h {
		quickSortAlg(res, 0, len(res)-1, flags, &CompareBySuffix{})
	}

	if flags.r {
		res = reverseString(res)
	}
	if flags.u {
		res = uniqueStrs(res)
	}
	return res
}

// возвращаем строки в обратном порядке
func reverseString(s []string) []string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

// возвращаем уникальные строки
func uniqueStrs(s []string) []string {
	strs := make(map[string]int)
	sOut := make([]string, 0, cap(s))
	for _, v := range s {
		if _, ok := strs[v]; !ok {
			sOut = append(sOut, v)
			strs[v]++
		}
	}
	return sOut
}

// проверка на отсортированность, флаг -c
func checkSort(s []string, flags flagsCmd) bool {
	if flags.k < 0 && !flags.n && !flags.M && !flags.h {
		for i := 0; i < len(s)-1; i++ {
			if s[i] > s[i+1] {
				return false
			}
		}
	}
	if flags.k > 0 && !flags.n && !flags.M && !flags.h {
		for i := 0; i < len(s)-1; i++ {
			cmp := &CompareByColumn{}
			if cmp.compare(s[i], s[i+1], flags) > 0 {
				return false
			}
		}
	}
	if flags.k < 0 && flags.n && !flags.M && !flags.h {
		for i := 0; i < len(s)-1; i++ {
			cmp := &CompareByInt{}
			if cmp.compare(s[i], s[i+1], flags) > 0 {
				return false
			}
		}
	}
	if flags.k < 0 && !flags.M && flags.h {
		for i := 0; i < len(s)-1; i++ {
			cmp := &CompareBySuffix{}
			if cmp.compare(s[i], s[i+1], flags) > 0 {
				return false
			}
		}
	}
	if flags.k < 0 && flags.M {
		for i := 0; i < len(s)-1; i++ {
			cmp := &CompareByMonth{}
			if cmp.compare(s[i], s[i+1], flags) > 0 {
				return false
			}
		}
	}
	return true
}

// в качестве сортировки используем quicksort, передаем компаратор в зависимости от флагов
func quickSortAlg(arr []string, lowIndex, highIndex int, fl flagsCmd, cmp Comparator) {
	if lowIndex >= highIndex {
		return
	}
	s2 := rand.NewSource(time.Now().Unix())
	r2 := rand.New(s2)
	pivotIndex := r2.Intn(highIndex-lowIndex) + lowIndex
	pivot := arr[pivotIndex]
	arr[pivotIndex], arr[highIndex] = arr[highIndex], arr[pivotIndex]

	leftP := lowIndex
	rightP := highIndex - 1
	for leftP < rightP {
		for cmp.compare(arr[leftP], pivot, fl) <= 0 && leftP < rightP {
			leftP++
		}
		for cmp.compare(arr[rightP], pivot, fl) >= 0 && leftP < rightP {
			rightP--
		}
		//свап
		arr[leftP], arr[rightP] = arr[rightP], arr[leftP]
	}

	if cmp.compare(arr[leftP], arr[highIndex], fl) >= 1 {
		arr[leftP], arr[highIndex] = arr[highIndex], arr[leftP]
	} else {
		leftP = highIndex
	}

	quickSortAlg(arr, lowIndex, leftP-1, fl, cmp)
	quickSortAlg(arr, leftP+1, highIndex, fl, cmp)

}

// компаратор, возвращает 1 если a > b; -1 если a < b; 0 если a=b
type Comparator interface {
	compare(a, b string, f flagsCmd) int
}

type CompareByDefault struct{}

func (c *CompareByDefault) compare(a, b string, f flagsCmd) int {
	a = strings.ToLower(a)
	b = strings.ToLower(b)
	if f.b {
		a = strings.TrimSpace(a)
		b = strings.TrimSpace(b)
	}
	if a > b {
		return 1
	} else if a < b {
		return -1
	}
	return 0
}

// сравнение по суффиксу
type CompareBySuffix struct{}

func (c *CompareBySuffix) compare(a, b string, f flagsCmd) int {
	a = strings.ToLower(a)
	b = strings.ToLower(b)
	if f.b {
		a = strings.TrimSpace(a)
		b = strings.TrimSpace(b)
	}
	if f.k > 0 {
		a = strings.Split(a, " ")[f.k-1]
		b = strings.Split(b, " ")[f.k-1]
	}

	if _, ok := Numsufx[string(a[len(a)-1])]; !ok {
		return -1
	}
	if _, ok := Numsufx[string(b[len(b)-1])]; !ok {
		return 1
	}
	v1 := Numsufx[string(a[len(a)-1])]
	v2 := Numsufx[string(b[len(b)-1])]
	v1fl, err := strconv.ParseFloat(a[:len(a)-1], 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error converting to float\n")
		return -1
	}
	v2fl, err := strconv.ParseFloat(b[:len(b)-1], 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error converting to float\n")
		return 1
	}
	v1fl = math.Pow(10, v1) * v1fl
	v2fl = math.Pow(10, v2) * v2fl
	if v1fl > v2fl {
		return 1
	} else if v1fl < v2fl {
		return -1
	}
	return 0
}

// компаратор по месяцу
type CompareByMonth struct{}

func (c *CompareByMonth) compare(a, b string, f flagsCmd) int {
	if f.b {
		a = strings.TrimSpace(a)
		b = strings.TrimSpace(b)
	}
	as := strings.Split(a, " ")
	bs := strings.Split(b, " ")
	var aMonth, bMonth byte
	for _, v := range as {
		if val, ok := Months[strings.ToLower(v)]; ok {
			aMonth = val
			break
		}
	}
	for _, v := range bs {
		if val, ok := Months[strings.ToLower(v)]; ok {
			bMonth = val
			break
		}
	}
	if aMonth > bMonth {
		return 1
	}
	if aMonth < bMonth {
		return -1
	}
	if aMonth == bMonth {
		return 0
	}
	return 0
}

// компаратор двух строк по колонке
type CompareByColumn struct{}

func (c *CompareByColumn) compare(a, b string, f flagsCmd) int {
	if f.b {
		a = strings.TrimSpace(a)
		b = strings.TrimSpace(b)
	}
	as := strings.Split(a, " ")[f.k-1]
	bs := strings.Split(b, " ")[f.k-1]

	//если по числовому значению
	if f.n {
		an, err := strconv.Atoi(as)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error converting to int\n")
			return -1
		}
		bn, err := strconv.Atoi(bs)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error converting to int\n")
			return 1
		}
		if an > bn {
			return 1
		}
		if an < bn {
			return -1
		}
		if an == bn {
			return 0
		}
	}

	if as > bs {
		return 1
	}
	if as < bs {
		return -1
	}
	if as == bs {
		return 0
	}
	return 0
}

// компаратор двух строк по числовым значениям
type CompareByInt struct{}

func (c *CompareByInt) compare(a, b string, f flagsCmd) int {
	//возвращает 1, если a > b поэлементно, -1 иначе, 0 если равны
	//конвертируем в массив интов
	if f.b {
		a = strings.TrimSpace(a)
		b = strings.TrimSpace(b)
	}
	as := strings.Split(a, " ")
	bs := strings.Split(b, " ")
	aInt := make([]int, 0)
	bInt := make([]int, 0)
	for _, v := range as {
		val, err := strconv.Atoi(v)
		if err != nil {
			continue
		}
		aInt = append(aInt, val)
	}
	for _, v := range bs {
		val, err := strconv.Atoi(v)
		if err != nil {
			continue
		}
		bInt = append(bInt, val)
	}

	if len(aInt) <= len(bInt) {
		for i := 0; i < len(aInt); i++ {
			if aInt[i] < bInt[i] {
				return -1
			}
			if aInt[i] > bInt[i] {
				return 1
			}
		}
		if len(aInt) == len(bInt) {
			return 0
		}
		return -1
	} else if len(bInt) < len(aInt) {
		for i := 0; i < len(bInt); i++ {
			if aInt[i] < bInt[i] {
				return -1
			}
			if aInt[i] > bInt[i] {
				return 1
			}
		}
		return 1
	}
	return 0
}

type flagsCmd struct {
	k int
	n bool
	r bool
	u bool
	M bool
	b bool
	c bool
	h bool
}

func flagsInit() *flagsCmd {
	f := new(flagsCmd)
	flag.IntVar(&f.k, "k", -1, "sort column")
	flag.BoolVar(&f.n, "n", false, "sort by numeric value")
	flag.BoolVar(&f.r, "r", false, "sort reverse")
	flag.BoolVar(&f.u, "u", false, "no repeats")
	flag.BoolVar(&f.M, "M", false, "sort by month name")
	flag.BoolVar(&f.b, "b", false, "ignore suffix space")
	flag.BoolVar(&f.c, "c", false, "check if already sorted")
	flag.BoolVar(&f.h, "h", false, "sort by numeric value with suffix")
	flag.Parse()
	return f
}

func readFile() []string {
	file, err := os.Open("./file.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	strs := make([]string, 0)
	for scanner.Scan() {
		strs = append(strs, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return strs
}

func writeFile(s *[]string) {
	f, err := os.Create("./outfile.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	for _, v := range *s {
		w.WriteString(v + "\n")
	}
	w.Flush()
}

// пример запуска go run . -k 2
func main() {
	flags := flagsInit()
	strs := readFile()
	for _, v := range strs {
		fmt.Println(v)
	}
	fmt.Println("-------SORTED-------")
	sorted := quickSort(strs, *flags)
	for _, v := range sorted {
		fmt.Println(v)
	}
	writeFile(&sorted)
}
