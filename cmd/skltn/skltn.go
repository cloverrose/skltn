package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/atotto/clipboard"
	"github.com/cloverrose/skltn"
	"github.com/cloverrose/skltn/templates"
)

const usage = `Usage:
1. Copy the method you want to test
2. Run 'skltn [-f] -t TemplateName'
3. Paste the skeleton that has been written to your clipboard`

func main() {
	flag.Usage = func() { fmt.Println(usage) }
	templateName := flag.String("t", "", "TemplateName")
	formatFlag := flag.Bool("f", false, "Format")
	flag.Parse()
	if *templateName == "" {
		errorOut(errors.New("invalid TemplateName"))
	}
	tplStr, err := templates.Get(*templateName)
	if err != nil {
		errorOut(err)
	}

	s, err := clipboard.ReadAll()
	if err != nil {
		errorOut(err)
	}
	mock, err := skltn.Generate([]byte(s), tplStr, *formatFlag)
	if err != nil {
		errorOut(err)
	}

	if err = clipboard.WriteAll(string(mock)); err != nil {
		errorOut(err)
	}
}

func errorOut(err error) {
	fmt.Fprintln(os.Stderr, err.Error())
	flag.Usage()
	os.Exit(1)
}
