Table of Contents
=================

  * [Table of Contents](#table-of-contents)
  * [fwd \- The little forwarder that could](#fwd---the-little-forwarder-that-could)
    * [About](#about)
    * [Use Cases](#use-cases)
      * [fwd ♥️ ngrok](#fwd-%EF%B8%8F-ngrok)
      * [Simple Use Case](#simple-use-case)
    * [Install](#install)
      * [Binary Releases](#binary-releases)
      * [Go Tool](#go-tool)
    * [Usage](#usage)
    * [Credits](#credits)

# fwd - The little forwarder that could
[![Travis](https://img.shields.io/travis/kintoandar/fwd.svg)](https://travis-ci.org/kintoandar/fwd)

## About
`fwd` is a network port forwarder written in golang.

It's cross platform, supports multiple architectures and it's dead simple to use.

Read all about it in this [article](https://blog.kintoandar.com/2016/08/fwd-the-little-forwarder-that-could.html).

## Use Cases
### fwd ♥️ ngrok
I must admit `ngrok` was an huge inspiration for `fwd`. If you don't know the tool you should definitely check out [this talk](https://www.youtube.com/watch?v=F_xNOVY96Ng) from [@inconshreveable](https://twitter.com/inconshreveable).

This tool combo (fwd + ngrok) allows some wicked mischief, like taking [firewall hole punching](https://en.wikipedia.org/wiki/Hole_punching_(networking)) to another level! And the setup is trivial.

`ngrok` allows to expose a local port on a public endpoint and `fwd` allows to connect a local port to a remote endpoint. You see where I'm heading with this... With both tools you can connect a public endpoint to a remote port as long as you have access to it.

Here's how it works:

```
                              +---------+                            +---------+
                        :9000 |         |            172.28.128.3:22 |         |
Internet +------------------> |   fwd   | +------------------------> |   ssh   |
tcp.ngrok.io:1234             |         | 172.28.128.1               |         |
                              +---------+                            +---------+
```

```
# get a public endpoint, ex: tcp.ngrok.io:1234
ngrok tcp 9000

# forward connections on :9000 to 172.28.128.3:22
fwd --from :9000 --to 172.28.128.3:22

# get a shell on 172.28.128.3 via a public endpoint
ssh tcp.ngrok.io -p 1234
```
_With great power comes great responsibility._ - Ben Parker

### Simple Use Case
Forwarding a local port to a remote port on a different network:

```
                             +---------+                             +---------+
           192.168.1.99:8000 |         |             172.28.128.3:80 |         |
curl +---------------------> |   fwd   | +-------------------------> |   web   |
                             |         | 172.28.128.1                |         |
                             +---------+                             +---------+
```

![demo](https://docs.google.com/uc?id=0B-SEc73VBiUwN0RheHVYQ3RlbW8)

## Install
Get the binaries or build it yourself.

### Binary Releases
Download prebuilt binaries for several platforms and architectures:

[![bintray](https://docs.google.com/uc?id=0B-SEc73VBiUwQ0NNLWRXdUN1M3c)](https://dl.bintray.com/kintoandar/fwd/)
[Bintray](https://dl.bintray.com/kintoandar/fwd/)

### Go Tool
```
go get github.com/kintoandar/fwd
go install github.com/kintoandar/fwd
```

## Usage
```
NAME:
   fwd - The little forwarder that could

USAGE:
   fwd --from localhost:2222 --to 192.168.1.254:22

VERSION:
   0.1.0

AUTHOR(S):
   Joel Bastos <kintoandar@gmail.com>

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --from value, -f value  source HOST:PORT (default: "127.0.0.1:8000") [$FWD_FROM]
   --to value, -t value    destination HOST:PORT [$FWD_TO]
   --list, -l              list local addresses
   --udp, -u               enable udp forwarding (tcp by default)
   --help, -h              show help
   --version, -v           print the version

COPYRIGHT:
   MIT License
```

## Credits
Made with ♥️ by [kintoandar](https://blog.kintoandar.com)
