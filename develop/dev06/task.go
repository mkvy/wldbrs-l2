package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

/*
=== Утилита cut ===

# Принимает STDIN, разбивает по разделителю (TAB) на колонки, выводит запрошенные

Поддержать флаги:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/
type Flags struct {
	fields    []int
	delimiter string
	sep       bool
}

func cutStrings(strs []string, f Flags) []string {
	//выходной слайс
	outStr := make([]string, 0, len(strs)/2)
	for _, v := range strs {
		//идем по входному слайсу, разделяем делимитером
		splitted := strings.Split(v, f.delimiter)
		//если условие что строка должна иметь разделитель, пропускаем слайсы разбитых строк из одного элемента
		if f.sep {
			if len(splitted) <= 1 {
				continue
			}
		}
		//собираем каждую строку
		strLine := make([]string, 0, len(splitted))
		//условие на добавление в выходную строку
		notEmpty := false
		for _, val := range f.fields {
			//идем по переданным параметрам fields, если выходим за количество разбитых подстрок, идем дальше
			if val-1 < len(splitted) {
				strLine = append(strLine, splitted[val-1])
				notEmpty = true
			}

		}
		if notEmpty {
			//собираем строку на выход и добавляем в выходной слайс
			outStrLine := strings.Join(strLine, f.delimiter)
			outStr = append(outStr, outStrLine)
		}
	}
	return outStr
}

// парсим параметры fields
func parseFlagString(s string) []int {
	strs := strings.Split(s, ",")
	ints := make([]int, 0, len(strs))
	for _, v := range strs {
		val, err := strconv.Atoi(v)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error converting string to int")
		}
		ints = append(ints, val)
	}
	return ints
}

func flagsInit() *Flags {
	f := &Flags{}
	stringInts := flag.String("f", "", "fields")
	flag.StringVar(&f.delimiter, "d", "\t", "delimiter")
	flag.BoolVar(&f.sep, "s", false, "only strings with delimiter")
	flag.Parse()
	f.fields = parseFlagString(*stringInts)
	return f
}

// читаем строки с stdin до сигнала eof и возвращаем слайс из считанных строк
func readStrings() []string {
	reader := bufio.NewReader(os.Stdin)
	strs := make([]string, 0)
	fmt.Println("Reading strings until EOF (ctrl+z / ctrl+d)")
	for {
		str, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		str = strings.TrimSpace(str)
		str = strings.TrimSuffix(str, "\n")
		strs = append(strs, str)
		//обрабатываем ctrl z в win / ctrl d linux
	}
	return strs
}

func printStrings(s []string) {
	fmt.Println("--------------")
	for _, v := range s {
		fmt.Println(v)
	}
}

func StartCut() {
	//читаем флаги
	flags := *flagsInit()
	//читаем строки
	strsToCut := readStrings()
	//обрезаем, печатаем
	out := cutStrings(strsToCut, flags)
	printStrings(out)
}

// пример запуска go run . -f 1,2 -d " " -s
func main() {
	StartCut()
}
