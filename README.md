# dockerjs

Docker.js is a javascript runtime for docker. It allows you to write javascript programs that orchestrate docker containers. Docker.js is also an interactive javascript shell that lets you interact with your Docker runtime. 

# Installation 

```
$ git clone https://github.com/jsr/dockerjs.git
$ cd dockerjs 
$ go get 
$ go build 
```

In order to run docker.js, you'll need a docker host available and your ENV setup so that the docker client API can reach it. 

```
$ ./dockerjs
Docker.js shell version 0.0.1
Type help() to get a list of builtin commands

``` 

You can create a new container 
```
> c = new Container("redis", "my-first-container") 
{
  "ID": "76223d0a50e53f4ca0159767aff0db12bac2922b5fde300f38c13f2e3864dd0e",
  "image": "redis",
  "name": "my-first-container"
}
``` 
You can register callbacks on container events 
``` 
> c.on('start', function(event) { print("hello world!"); })
```
Now we'll start the container. You should see "hello world!" printed when the start event hits. 
```
> c.run()
>
hello world!
>
```
You can get a list of all of the containers in the runtime
```
> Container.all()
[
  {
    "ID": "76223d0a50e53f4ca0159767aff0db12bac2922b5fde300f38c13f2e3864dd0e",
    "image": "redis",
    "name": "/my-first-container"
  }
]
>
``` 

You can also search for a specific container by ID (or prefix of ID) 
```
> Container.findByID("76223d0a50e5")
{
  "ID": "76223d0a50e53f4ca0159767aff0db12bac2922b5fde300f38c13f2e3864dd0e",
  "image": "redis",
  "name": "/my-first-container"
}
```


# High level API 


function|description
---|---
new Container(image,name) | Creates a new container 
Container.all() | returns a list of all running containers 
Container.findByID(id) | returns container matching the specified id 
Container.match(fn) | fn is a function which is evaluated against each running container. matched containers are returned 
Container.prototype.run() | starts the container 
Container.prototype.on(event, fn) | calls fn when the specified event is raised by the container 


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
