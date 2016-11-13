package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/chzyer/readline"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/robertkrimen/otto"

	"io"
	"io/ioutil"
)

func main() {
	js := otto.New()
	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println("Couldn't connect to docker: ", err)
		return
	}

	js.Set("require", func(call otto.FunctionCall) otto.Value {
		filename, _ := call.Argument(0).ToString()
		file, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Println("Couldn't read file: ", filename)
			val, _ := js.ToValue(nil)
			return val
		}
		result, err := js.Run(file)
		if err != nil {
			fmt.Println("Couldn't execute file: ", err)
			val, _ := js.ToValue(nil)
			return val

		}
		val, _ := js.ToValue(result)
		return val
	})

	js.Set("print", func(call otto.FunctionCall) otto.Value {
		msg, _ := call.Argument(0).ToString()
		fmt.Println(msg)
		val, _ := js.ToValue(nil)
		return val
	})

	js.Set("containers", func(call otto.FunctionCall) otto.Value {
		containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
		if err != nil {
			fmt.Println("Couldn't find containers: ", err)
			val, _ := js.ToValue(nil)
			return val
		}

		var results []otto.Value

		for _, container := range containers {
			val, _ := js.ToValue(container)
			results = append(results, val)
		}
		val, _ := js.ToValue(results)
		return val
	})

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

		value, err := js.Run(input)
		if err != nil {
			fmt.Println("Error running command: ", err)
		}
		printOttoValue(value)
	}
}

func printOttoValue(value otto.Value) {
	if value.IsDefined() {
		if value.Class() == "GoArray" {
			arr, _ := value.Export()
			for _, elt := range arr.([]otto.Value) {
				obj, _ := elt.Export()
				jsn, _ := json.Marshal(obj)
				fmt.Println(" ", string(jsn))
			}
		} else {
			obj, _ := value.Export()
			jsn, _ := json.Marshal(obj)
			fmt.Println(string(jsn))
		}
	}
}
