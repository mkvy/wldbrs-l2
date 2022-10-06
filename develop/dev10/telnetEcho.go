package main

import (
	"github.com/reiver/go-telnet"
	"time"
)

func GoTelnet() {
	var handler = telnet.EchoHandler
	time.AfterFunc(15*time.Second, func() { return })
	err := telnet.ListenAndServe(":5555", handler)
	if nil != err {
		panic(err)
	}
}
