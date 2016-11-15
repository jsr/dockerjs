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
    "Command": "docker-entrypoint.sh redis-server",
    "Created": 1479248067,
    "HostConfig": {
      "NetworkMode": "default"
    },
    "ID": "76223d0a50e53f4ca0159767aff0db12bac2922b5fde300f38c13f2e3864dd0e",
    "Image": "redis",
    "ImageID": "sha256:5f515359c7f871387c7dbb9bcf5c56b971d6b48ec9fec9e68d9369f4673db218",
    "Labels": {},
    "Mounts": [
      {
        "Destination": "/data",
        "Driver": "local",
        "Mode": "",
        "Name": "1995df81ed898962dc2b05812953386404d2553151a7935a9d4cdb81971b1313",
        "Propagation": "",
        "RW": true,
        "Source": "/var/lib/docker/volumes/1995df81ed898962dc2b05812953386404d2553151a7935a9d4cdb81971b1313/_data",
        "Type": ""
      }
    ],
    "Names": [
      "/my-first-container"
    ],
    "NetworkSettings": {
      "Networks": {
        "bridge": {
          "Aliases": [],
          "EndpointID": "7c5a39e54432ed98ed994683800990db522085ca294cd5b0e9515b5bd2f906f6",
          "Gateway": "172.17.0.1",
          "GlobalIPv6Address": "",
          "GlobalIPv6PrefixLen": 0,
          "IPAddress": "172.17.0.2",
          "IPPrefixLen": 16,
          "IPv6Gateway": "",
          "Links": [],
          "MacAddress": "02:42:ac:11:00:02",
          "NetworkID": "fa844df75e2f649827da093c3870643f6f7e18b9f60efe7d3ec09f5215fe7ca8"
        }
      }
    },
    "Ports": [
      {
        "IP": "",
        "PrivatePort": 6379,
        "PublicPort": 0,
        "Type": "tcp"
      }
    ],
    "SizeRootFs": 0,
    "SizeRw": 0,
    "State": "running",
    "Status": "Up 7 seconds"
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
