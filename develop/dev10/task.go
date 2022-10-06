package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

/*
=== Утилита telnet ===

Реализовать примитивный telnet клиент:
Примеры вызовов:
go-telnet --timeout=10s host port go-telnet mysite.ru 8080 go-telnet --timeout=3s 1.1.1.1 123

Программа должна подключаться к указанному хосту (ip или доменное имя) и порту по протоколу TCP.
После подключения STDIN программы должен записываться в сокет, а данные полученные и сокета должны выводиться в STDOUT
Опционально в программу можно передать таймаут на подключение к серверу (через аргумент --timeout, по умолчанию 10s).

При нажатии Ctrl+D программа должна закрывать сокет и завершаться. Если сокет закрывается со стороны сервера, программа должна также завершаться.
При подключении к несуществующему сервер, программа должна завершаться через timeout.
*/

type TelnetClient struct {
	timeout time.Duration
	host    string
	port    string
}

func NewTelnetClient(timeout time.Duration, host, port string) *TelnetClient {
	//таймаут в секундах
	return &TelnetClient{
		timeout: timeout * time.Second,
		host:    host,
		port:    port,
	}
}

func (t *TelnetClient) Run() {
	//устанавливаем соединение с таймаутом
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(t.host, t.port), t.timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	for {
		fmt.Print(">> ")
		//читаем данные c stdin
		reader := bufio.NewReader(os.Stdin)
		clientReq, err := reader.ReadString('\n')
		if err != nil {
			//если был сигнал eof ctrl+z
			if err == io.EOF {
				fmt.Fprintln(os.Stdout, "Closing connection")
				return
			}
			fmt.Fprintln(os.Stderr, "error: ", err)
			return
		}
		//пишем в соединение
		_, err = fmt.Fprint(conn, clientReq)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing message: ", err)
		}
		//читаем данные с соединения
		servReader := bufio.NewReader(conn)
		serverRes, err := servReader.ReadString('\n')
		if err != nil {
			//если был сигнал eof ctrl+z
			if err == io.EOF {
				fmt.Fprintln(os.Stdout, "Server closed the connection")
				return
			}
			fmt.Fprintln(os.Stderr, "error: ", err)
			return
		}
		//пишем ответ
		fmt.Fprint(os.Stdout, "->:")
		fmt.Fprint(os.Stdout, serverRes)
	}
}

// пример запуска go run . -timeout 15 localhost 5555
func main() {
	timeout := flag.Int("timeout", 10, "connection timeout")
	flag.Parse()
	t := *timeout
	telnet := NewTelnetClient(time.Duration(t), flag.Arg(0), flag.Arg(1))
	go GoTelnet()
	telnet.Run()
}
