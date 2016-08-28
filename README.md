# fwd - The little forwarder that could

## Introduction
`fwd` is a TCP/UDP forwarder written in golang

## Install
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

# Credits
Made with ♥️ by [kintoandar](https://blog.kintoandar.com)
