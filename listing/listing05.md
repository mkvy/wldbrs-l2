Что выведет программа? Объяснить вывод программы.

```go
package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	{
		// do something
	}
	return nil
}

func main() {
	var err error
	err = test()
	if err != nil {
		println("error")
		return
	}
	println("ok")
}
```

Ответ:
```
Выведет error.
Аналогично вопросу из листинга 03.
Интерфейс имеет два поля - dynamicValue, dynamicType.
dynamicType в данном случае не равен nil, поэтому проверка на == nil не прошла и вывелось error
```
