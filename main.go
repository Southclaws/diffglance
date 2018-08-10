package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/nsf/termbox-go"
	"github.com/sergi/go-diff/diffmatchpatch"
)

func main() {
	err := run()
	if err != nil {
		fmt.Println(err)
	}
}

func run() (err error) {
	var (
		a *os.File
		b *os.File
	)

	flag.Parse()
	switch flag.NArg() {
	case 2:
		a, err = os.Open(flag.Arg(0))
		if err != nil {
			fmt.Println("failed to open input file:", err)
		}
		b, err = os.Open(flag.Arg(1))
		if err != nil {
			fmt.Println("failed to open input file:", err)
		}
	default:
		fmt.Printf("input must be from stdin or file\n")
		os.Exit(1)
	}

	abytes, err := ioutil.ReadAll(a)
	if err != nil {
		return
	}

	bbytes, err := ioutil.ReadAll(b)
	if err != nil {
		return
	}

	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(string(abytes), string(bbytes), false)

	err = termbox.Init()
	if err != nil {
		return
	}
	width, termHeight := termbox.Size()
	if termHeight == 0 {
		err = errors.New("Could not get terminal height")
		termbox.Close()
		return
	}
	termbox.Close()

	contentHeight := 0
	for _, diff := range diffs {
		contentHeight += len(strings.Split(diff.Text, "\n"))
	}

	for _, diff := range diffs {
		text := diff.Text
		textLines := len(strings.Split(text, "\n"))
		blockHeight := (float64(textLines) / float64(contentHeight)) * float64(termHeight)

		var colour color.Attribute

		switch diff.Type {
		case diffmatchpatch.DiffInsert:
			colour = color.FgGreen
		case diffmatchpatch.DiffDelete:
			colour = color.FgRed
		case diffmatchpatch.DiffEqual:
			colour = color.FgBlack
		}

		for index := 0; index < int(blockHeight); index++ {
			color.Set(colour)
			for x := 0; x < width; x++ {
				fmt.Print("#")
			}
			fmt.Print("\n")
		}
	}

	return nil
}
