Table of Contents
=================

  * [Table of Contents](#table-of-contents)
  * [fwd \- The little forwarder that could](#fwd---the-little-forwarder-that-could)
    * [About](#about)
    * [Demo](#demo)
    * [Install](#install)
      * [Binary releases](#binary-releases)
      * [Go tool](#go-tool)
    * [Usage](#usage)
    * [Credits](#credits)

# fwd - The little forwarder that could
[![Travis](https://img.shields.io/travis/kintoandar/fwd.svg)](https://travis-ci.org/kintoandar/fwd)

## About
`fwd` is a network port forwarder written in golang

## Demo
```
                             +---------+                             +---------+
           192.168.1.99:8000 |         |             172.28.128.3:80 |         |
curl +---------------------> |   fwd   | +-------------------------> |   web   |
                             |         | 172.28.128.1                |         |
                             +---------+                             +---------+
```

![demo](https://docs.google.com/uc?id=0B-SEc73VBiUwN0RheHVYQ3RlbW8)

## Install
Get the binaries or build it yourself

### Binary releases
Download prebuilt binaries for several platforms and architectures:

[![bintray](https://docs.google.com/uc?id=0B-SEc73VBiUwQ0NNLWRXdUN1M3c)](https://dl.bintray.com/kintoandar/fwd/) [Bintray](https://dl.bintray.com/kintoandar/fwd/)

### Go tool
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
