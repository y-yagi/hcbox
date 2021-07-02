package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"jaytaylor.com/html2text"
)

const cmd = "danwa"

var (
	flags *flag.FlagSet
)

func usage() {
	fmt.Fprintf(os.Stdout, "Usage: %s [OPTIONS] URL\n", cmd)
	fmt.Fprintln(os.Stdout, "OPTIONS:")
	flags.PrintDefaults()
}

func setFlags() {
	flags = flag.NewFlagSet(cmd, flag.ExitOnError)
	flags.Usage = usage
}

func main() {
	setFlags()
	err := flags.Parse(os.Args[1:])
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(2)
	}

	if len(flags.Args()) != 1 {
		usage()
		os.Exit(1)
	}

	baseurl := flags.Args()[0]

	l, err := readline.NewEx(&readline.Config{
		Prompt:          ">> ",
		InterruptPrompt: "^C",
	})
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(2)
	}

	defer l.Close()

	for {
		url := baseurl
		input, err := l.Readline()
		if err != nil {
			return
		}

		if len(input) == 0 {
			continue
		}

		methodAndPath := strings.Split(input, " ")
		if len(methodAndPath) != 1 {
			url = url + methodAndPath[1]
		}

		req, err := http.NewRequest(methodAndPath[0], url, nil)
		if err != nil {
			fmt.Printf("%v\n", err)
			continue
		}

		c := &http.Client{}
		res, err := c.Do(req)
		if err != nil {
			fmt.Printf("%v\n", err)
			continue
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("%v\n", err)
			continue
		}

		defer res.Body.Close()

		d := color.New(color.FgGreen, color.Bold)
		d.Println("Headers")
		for name, values := range res.Header {
			color.Cyan("  %s:  \n", name)
			for _, value := range values {
				fmt.Printf("    %s\n", value)
			}
		}

		text, err := html2text.FromString(string(body), html2text.Options{PrettyTables: true})
		if err != nil {
			fmt.Printf("%v\n", err)
			continue
		}

		if len(text) != 0 {
			d.Println("\n\nBody")
			fmt.Println(text)
		}
	}
}
