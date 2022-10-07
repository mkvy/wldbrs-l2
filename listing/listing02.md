Что выведет программа? Объяснить вывод программы. Объяснить как работают defer’ы и их порядок вызовов.

```go
package main

import (
	"fmt"
)


func test() (x int) {
	defer func() {
		x++
	}()
	x = 1
	return
}


func anotherTest() int {
	var x int
	defer func() {
		x++
	}()
	x = 1
	return x
}


func main() {
	fmt.Println(test())
	fmt.Println(anotherTest())
}
```

Ответ:
```
test() выведет 2, anotherTest() выведет 1
test() имеет именованный выходной параметер
Анонимные defer функции имеют доступ к именованным выходным параметрам.
Во втором случае параметер не именованный 
и анонимная функция работает с копией переменной, поэтому ни на что не влияет
defer выполняется сразу после return

```
