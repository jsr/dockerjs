package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/chzyer/readline"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
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
			return otto.NullValue()
		}
		result, err := js.Run(file)
		if err != nil {
			fmt.Println("Couldn't execute file: ", err)
			return otto.NullValue()

		}
		val, _ := js.ToValue(result)
		return val
	})

	js.Set("print", func(call otto.FunctionCall) otto.Value {
		msg, _ := call.Argument(0).ToString()
		fmt.Println(msg)
		return otto.NullValue()
	})

	js.Set("containers", func(call otto.FunctionCall) otto.Value {
		containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
		if err != nil {
			fmt.Println("Couldn't find containers: ", err)
			return otto.NullValue()
		}

		var results []otto.Value

		for _, container := range containers {
			val, _ := js.ToValue(container)
			results = append(results, val)
		}
		val, _ := js.ToValue(results)
		return val
	})

	js.Set("run", func(call otto.FunctionCall) otto.Value {
		containerConfigRaw, err := call.Argument(0).Export()
		if err != nil {
			fmt.Println("Must specify container.Config")
			return otto.NullValue()
		}
		containerConfigJson, _ := json.Marshal(containerConfigRaw)
		containerConfig := new(container.Config)
		json.Unmarshal(containerConfigJson, containerConfig)

		hostConfigRaw, err := call.Argument(1).Export()
		if err != nil {
			fmt.Println("Must specifh container.HostConfig")
			return otto.NullValue()
		}

		hostConfigJson, _ := json.Marshal(hostConfigRaw)
		hostConfig := new(container.HostConfig)
		json.Unmarshal(hostConfigJson, hostConfig)

		networkConfigRaw, err := call.Argument(2).Export()
		if err != nil {
			fmt.Println("Must specify network.NetworkingConfig")
			return otto.NullValue()
		}

		networkConfigJson, _ := json.Marshal(networkConfigRaw)
		networkConfig := new(network.NetworkingConfig)
		json.Unmarshal(networkConfigJson, networkConfig)

		containerName, err := call.Argument(3).Export()
		if err != nil {
			fmt.Println("Must specify name for container")
			return otto.NullValue()
		}

		if len(containerConfig.Image) == 0 {
			fmt.Println("Must specify image of container to run")
			return otto.NullValue()
		}

		pull_result, err := cli.ImagePull(context.Background(), containerConfig.Image, types.ImagePullOptions{})
		if err != nil {
			fmt.Println("Cloudn't pull image:", err)
			return otto.NullValue()
		}
		body, err := ioutil.ReadAll(pull_result)
		if err != nil {
			fmt.Println("Couldn't pull image after read:", err)
			return otto.NullValue()
		}
		if bytes.Contains(body, []byte("errorDetail")) {
			fmt.Println("Error loading imag")
			return otto.NullValue()
		}

		create_result, err := cli.ContainerCreate(context.Background(), containerConfig, hostConfig, networkConfig, containerName.(string))

		if err != nil {
			fmt.Println("Failed creating container:", err)
			return otto.NullValue()
		}

		err = cli.ContainerStart(context.Background(), create_result.ID, types.ContainerStartOptions{})

		if err != nil {
			fmt.Println("Failed starting container:", err)
			return otto.NullValue()
		}

		val, _ := js.ToValue(create_result)
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
				var out bytes.Buffer
				json.Indent(&out, jsn, "", "\t")
				fmt.Println(out.String())
			}
		} else {
			obj, _ := value.Export()
			jsn, _ := json.Marshal(obj)
			var out bytes.Buffer
			json.Indent(&out, jsn, "", "\t")
			fmt.Println(out.String())
		}
	}
}
