package vm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	docker "github.com/docker/docker/client"
	"github.com/robertkrimen/otto"
)

type vm struct {
	otto   *otto.Otto
	docker *docker.Client
}

func New() *vm {
	vm := new(vm)
	docker, err := docker.NewEnvClient()
	if err != nil {
		panic("Couldn't connect to docker host")
	}
	vm.docker = docker

	vm.otto = otto.New()
	vm.installExtensions()
	return vm
}

func (vm *vm) Evaluate(input string) (output string, err error) {
	out, err := vm.otto.Run(input)
	if err != nil {
		return "", err
	}
	if !out.IsNull() {
		json, _ := vm.otto.Call("JSON.stringify", nil, out, nil, "  ")
		return json.String(), nil
	} else {
		return "", nil
	}
}

func (vm *vm) installExtensions() {
	vm.otto.Set("sleep", func(call otto.FunctionCall) otto.Value {
		millis, _ := call.Argument(0).ToInteger()
		time.Sleep(time.Duration(millis) * time.Millisecond)
		return otto.NullValue()
	})

	vm.otto.Set("require", func(call otto.FunctionCall) otto.Value {
		filename, _ := call.Argument(0).ToString()
		file, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Println("Couldn't read file: ", filename)
			return otto.NullValue()
		}
		result, err := vm.otto.Run(file)
		if err != nil {
			fmt.Println("Couldn't execute file: ", err)
			return otto.NullValue()

		}
		val, _ := vm.otto.ToValue(result)
		return val
	})

	vm.otto.Set("print", func(call otto.FunctionCall) otto.Value {
		msg, _ := call.Argument(0).ToString()
		fmt.Println(msg)
		return otto.NullValue()
	})

	vm.otto.Set("containers", func(call otto.FunctionCall) otto.Value {
		containers, err := vm.docker.ContainerList(context.Background(), types.ContainerListOptions{})
		if err != nil {
			fmt.Println("Couldn't find containers: ", err)
			return otto.NullValue()
		}

		var results []otto.Value

		for _, container := range containers {
			val, _ := vm.otto.ToValue(container)
			results = append(results, val)
		}
		val, _ := vm.otto.ToValue(results)
		return val
	})

	vm.otto.Set("run", func(call otto.FunctionCall) otto.Value {
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

		pull_result, err := vm.docker.ImagePull(context.Background(), containerConfig.Image, types.ImagePullOptions{})
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

		create_result, err := vm.docker.ContainerCreate(context.Background(), containerConfig, hostConfig, networkConfig, containerName.(string))

		if err != nil {
			fmt.Println("Failed creating container:", err)
			return otto.NullValue()
		}

		err = vm.docker.ContainerStart(context.Background(), create_result.ID, types.ContainerStartOptions{})

		if err != nil {
			fmt.Println("Failed starting container:", err)
			return otto.NullValue()
		}

		val, _ := vm.otto.ToValue(create_result)
		return val
	})
}
