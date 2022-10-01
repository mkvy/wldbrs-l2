package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

/*
=== Задача на распаковку ===

Создать Go функцию, осуществляющую примитивную распаковку строки, содержащую повторяющиеся символы / руны, например:
	- "a4bc2d5e" => "aaaabccddddde"
	- "abcd" => "abcd"
	- "45" => "" (некорректная строка)
	- "" => ""
Дополнительное задание: поддержка escape - последовательностей
	- qwe\4\5 => qwe45 (*)
	- qwe\45 => qwe44444 (*)
	- qwe\\5 => qwe\\\\\ (*)

В случае если была передана некорректная строка функция должна возвращать ошибку. Написать unit-тесты.

Функция должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func isSymbol(ch rune) bool {
	return !isDigit(ch) && string(ch) != `\`
}

func isDigit(ch rune) bool {
	return unicode.IsDigit(ch)
}

func writeMultipleRunes(sb *strings.Builder, r rune, count int) strings.Builder {
	count -= 1
	for i := 0; i < count; i++ {
		sb.WriteRune(r)
	}
	return *sb
}

func UnpackStr(src string) (string, error) {
	//проверям строку на корректность
	if src == "" {
		return "", errors.New("empty string")
	}
	rs := []rune(src)
	if isDigit(rs[0]) {
		return "", errors.New("wrong string")
	}
	for i := 1; i < len(rs)-1; i++ {
		if isDigit(rs[i]) && isDigit(rs[i+1]) && string(rs[i-1]) != `\` {
			return "", errors.New("wrong string")
		}
	}
	var sb strings.Builder
	var prevR rune
	_ = prevR
	var escCount int = 0
	for _, r := range src {
		//если символ и не было эскейп символов
		if isSymbol(r) && escCount == 0 {
			//записываем в буфер
			sb.WriteRune(r)
			prevR = r
		} else if isDigit(r) && escCount == 0 {
			//если цифра, пишем предыдущую руну в буфер num раз
			num, err := strconv.Atoi(string(r))
			if err != nil {
				return "", errors.New("error converting to num")
			}
			sb = writeMultipleRunes(&sb, prevR, num)
		} else if string(r) == `\` && escCount == 0 {
			//если встретили esc символ
			escCount++
		} else if escCount > 0 {
			//обрабатываем символ за слешем
			escCount = 0
			prevR = r
			sb.WriteRune(r)
		}

	}
	return sb.String(), nil
}

func main() {
	str, err := UnpackStr(`a4bc2d5e`)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error occured: %v \n", err)
	}
	fmt.Println(str)
	str, err = UnpackStr(`45`)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error occured: %v \n", err)
	}
	fmt.Println(str)

	str, err = UnpackStr(`qwe\4\5`)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error occured: %v \n", err)
	}
	fmt.Println(str)

	str, err = UnpackStr(`qwe\\5`)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error occured: %v \n", err)
	}
	fmt.Println(str)
	str, err = UnpackStr(`qwe\\a`)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error occured: %v \n", err)
	}
	fmt.Println(str)

}
