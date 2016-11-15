package main

import (
	"fmt"
	"github.com/chzyer/readline"
	"github.com/jsr/dockerjs/vm"
	"io"
)

func main() {
	vm := vm.New()

	l, err := readline.NewEx(&readline.Config{
		Prompt:          "> ",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		fmt.Println("Panic setting up readline. ", err)
		return
	}
	defer l.Close()

	fmt.Println("Docker.js shell version 0.0.1")
	fmt.Println("Type help() to get a list of builtin commands")

	for {
		input, err := l.Readline()
		if err == io.EOF {
			break
		} else if err == readline.ErrInterrupt {
			if len(input) == 0 {
				break
			} else {
				continue
			}
		}
		if err != nil {
			fmt.Println("Error reading input: ", err)
			return
		}

		out, err := vm.Evaluate(input)

		if err != nil {
			fmt.Println("Error running command: ", err)
		}

		if len(out) > 0 {
			fmt.Println(out)
		}
	}
}
