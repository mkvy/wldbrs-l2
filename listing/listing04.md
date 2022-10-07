Что выведет программа? Объяснить вывод программы.

```go
package main

func main() {
	ch := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
		}
	}()

	for n := range ch {
		println(n)
	}
}
```

Ответ:
```
Значения от 0 до 9 и fatal error deadlock.
Range ожидает закрытия канала, этого не происходит, 
соответственно блокируется после окончания работы горутины

```
