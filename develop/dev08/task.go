package main

import (
	"bufio"
	"fmt"
	ps "github.com/mitchellh/go-ps"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

/*
Необходимо реализовать свой собственный UNIX-шелл-утилиту с поддержкой ряда простейших команд:

- cd <args> - смена директории (в качестве аргумента могут быть то-то и то)
- pwd - показать путь до текущего каталога
- echo <args> - вывод аргумента в STDOUT
- kill <args> - "убить" процесс, переданный в качесте аргумента (пример: такой-то пример)
- ps - выводит общую информацию по запущенным процессам в формате *такой-то формат*

Так же требуется поддерживать функционал fork/exec-команд
*/

func echo(arg string) {
	arg = strings.Replace(arg, "echo ", "", 1)
	fmt.Fprintln(os.Stdout, arg)
}

func cd(args []string) {
	err := os.Chdir(args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error with directory\n")
	}
}

func pwd() {
	val, _ := os.Getwd()
	fmt.Fprintln(os.Stdout, val)
}

func psGet() {
	prs, err := ps.Processes()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error searching process\n")
	}
	fmt.Fprintln(os.Stdout, "# pid ppid executable")
	for i, v := range prs {
		fmt.Fprintln(os.Stdout, i, v.Pid(), v.PPid(), v.Executable())
	}
}

func kill(args []string) {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "no pid entered!")
	}
	val, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error with pid format")
	}
	p, err := os.FindProcess(val)
	err = p.Kill()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error killing process")
	}
}

func execCmd(args []string) {
	cmd := exec.Command(args[1], args[2:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running exec\n")
	}
}

func parseStr(str string) {
	arg := strings.Split(str, " ")
	cmd := strings.ToLower(arg[0])
	//вызываем функцию в зависимости от команды
	switch cmd {
	case "echo":
		echo(str)
	case "cd":
		cd(arg)
	case "pwd":
		pwd()
	case "ps":
		psGet()
	case "kill":
		kill(arg)
	case "exec":
		execCmd(arg)
	}
}

func StartShell() {
	scan := bufio.NewScanner(os.Stdin)
	//в цикле считываем команды
	for {
		fmt.Fprint(os.Stdout, ">")
		scan.Scan()
		str := scan.Text()
		if str == "quit" {
			break
		} else {
			parseStr(str)
		}
	}
}

func main() {
	StartShell()
}
