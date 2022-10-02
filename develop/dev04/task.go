package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

/*
=== Поиск анаграмм по словарю ===

Напишите функцию поиска всех множеств анаграмм по словарю.
Например:
'пятак', 'пятка' и 'тяпка' - принадлежат одному множеству,
'листок', 'слиток' и 'столик' - другому.

Входные данные для функции: ссылка на массив - каждый элемент которого - слово на русском языке в кодировке utf8.
Выходные данные: Ссылка на мапу множеств анаграмм.
Ключ - первое встретившееся в словаре слово из множества
Значение - ссылка на массив, каждый элемент которого, слово из множества. Массив должен быть отсортирован по возрастанию.
Множества из одного элемента не должны попасть в результат.
Все слова должны быть приведены к нижнему регистру.
В результате каждое слово должно встречаться только один раз.

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

//проверка: является ли word2 анаграмой word1
func isAnagram(word1, word2 string) bool {
	if len(word2) != len(word1) {
		return false
	}
	//количество вхождений в строку
	chars := make(map[rune]int)
	for _, v := range word1 {
		chars[v]++
	}
	for _, v := range word2 {
		//если не найдена нужная буква, false
		if _, ok := chars[v]; !ok {
			return false
		}
		//когда кончились буквы, удаляем из множества
		if chars[v]-1 == 0 {
			delete(chars, v)
		}
		chars[v]--
	}
	return true
}

//быстрая сортировка, для сортировки слайса на выходе
func quickSort(arr []string) []string {
	res := make([]string, len(arr))
	copy(res, arr)
	quickSortAlg(res, 0, len(res)-1)
	return res
}

func quickSortAlg(arr []string, lowIndex, highIndex int) {
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
		for arr[leftP] <= pivot && leftP < rightP {
			leftP++
		}
		for arr[rightP] >= pivot && leftP < rightP {
			rightP--
		}
		//свап
		arr[leftP], arr[rightP] = arr[rightP], arr[leftP]
	}

	if arr[leftP] > arr[highIndex] {
		arr[leftP], arr[highIndex] = arr[highIndex], arr[leftP]
	} else {
		leftP = highIndex
	}

	quickSortAlg(arr, lowIndex, leftP-1)
	quickSortAlg(arr, leftP+1, highIndex)

}

//возвращает слайс состоящий только из уникальных строк
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

func GetAnagramSet(strs []string) map[string][]string {
	//множество
	vocab := make(map[string][]string)
	for _, v := range strs {
		notFound := true
		for k, _ := range vocab {
			//если для анаграмы ключ уже существует, добавляем в значение-слайс мапы
			if isAnagram(strings.ToLower(k), strings.ToLower(v)) {
				vocab[k] = append(vocab[k], strings.ToLower(v))
				notFound = false
			}
		}
		//если впервые, добавляем в мапу
		if notFound {
			vocab[strings.ToLower(v)] = append(vocab[strings.ToLower(v)], strings.ToLower(v))
		}
	}
	for k, v := range vocab {
		//удаляем одноэлементные
		if len(v) <= 1 {
			delete(vocab, k)
		}
		//убираем дубли и сортируем
		vocab[k] = uniqueStrs(v)
		vocab[k] = quickSort(vocab[k])
	}
	return vocab
}

func main() {
	strs := []string{"Столик", "СЛИТОК", "пятак", "ятпка", "СТОлик", "тяпка", "тяпка", "пятак"}
	out := GetAnagramSet(strs)
	fmt.Println(out)
}
