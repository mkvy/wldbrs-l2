package main

import (
	"flag"
	"fmt"
	"os"
)

/*
=== Утилита wget ===

Реализовать утилиту wget с возможностью скачивать сайты целиком

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type Flags struct {
	lDepth int
	url    string
}

func wget(f Flags) {
	if f.lDepth < 1 {
		fmt.Fprintln(os.Stderr, "wrong depth")
		return
	}
	wd, _ := os.Getwd()
	err := os.Mkdir(wd+`\`+f.url, os.ModeDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error creating directory")
		return
	}
}

func flagsInit() *Flags {
	f := &Flags{}
	flag.StringVar(&f.url, "l", "\t", "url")
	flag.IntVar(&f.lDepth, "f", 1, "depth of download")
	return f
}

func main() {
}
