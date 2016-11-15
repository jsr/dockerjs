# dockerjs

Docker.js is a javascript runtime for docker. It allows you to write javascript programs that orchestrate docker containers. Docker.js is also an interactive javascript shell that lets you interact with your Docker runtime. 

# Installation 

```
$ git clone https://github.com/jsr/dockerjs.git
$ go build 
$ ./dockerjs
Docker.js shell version 0.0.1
Type help() to get a list of builtin commands
> print("hello!") 
hello!
>
```
# High level API 


function|description
---|---
new Container(image,name) | Creates a new container 
Container.all() | returns a list of all running containers 
Container.findByID(id) | returns container matching the specified id 
Container.match(fn) | fn is a function which is evaluated against each running container. matched containers are returned 
<container>.run() | starts the container 
<container>.on(event, fn) | calls fn when the specified event is raised by the container 


# Low level API 
Docker.js exposes some low level API's that let you access the docker runtime. 

function|description
---|---
list() -> []container |Returns a list of all running containers 
create(config, hostConfig, networkingConfig, name ) -> container | Creates a new container. Will pull image if it's not there 
run(id) | Starts the specified container 
require(filename) | load a .js file and evaluate it 
print(string) | prints to stdout 
sleep(millis) | sleeps for miillis milliseconds 
listen(id, event, callback) | registers a callback function that will be invoked when the specified event is emitted by the specified container id 
