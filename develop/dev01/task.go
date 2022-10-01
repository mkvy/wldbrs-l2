package main

import (
	"fmt"
	"github.com/beevik/ntp"
	"os"
	"time"
)

/*
=== Базовая задача ===

Создать программу печатающую точное время с использованием NTP библиотеки.Инициализировать как go module.
Использовать библиотеку https://github.com/beevik/ntp.
Написать программу печатающую текущее время / точное время с использованием этой библиотеки.

Программа должна быть оформлена с использованием как go module.
Программа должна корректно обрабатывать ошибки библиотеки: распечатывать их в STDERR и возвращать ненулевой код выхода в OS.
Программа должна проходить проверки go vet и golint.
*/
const ntpServer = "ntp5.stratum2.ru"

func GetTime() (time.Time, error) {
	//получаем время через метод Time библиотеки ntp
	return ntp.Time(ntpServer)
}

func main() {
	time, err := GetTime()
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error has occured %s", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "Current time is %v\n", time.Format("15.04.05"))
}
