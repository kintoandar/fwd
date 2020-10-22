package main

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/urfave/cli"
	"io"
	"log"
	"net"
	"os"
	"runtime"
)

func getLocalAddrs() ([]net.IP, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}
	var list []net.IP
	for _, addr := range addrs {
		v := addr.(*net.IPNet)
		if v.IP.To4() != nil {
			list = append(list, v.IP)
		}
	}
	return list, nil
}

func fwd(src net.Conn, remote string, proto string) {
	dst, err := net.Dial(proto, remote)
	errHandler(err)
	go func() {
		_, err = io.Copy(src, dst)
		errPrinter(err)
	}()
	go func() {
		_, err = io.Copy(dst, src)
		errPrinter(err)
	}()
}

func errHandler(err error) {
	if err != nil {
		color.Set(color.FgRed)
		fmt.Fprintf(os.Stderr, "[Error] %s\n", err.Error())
		color.Unset()
		os.Exit(1)
	}
}

// TODO: merge error handling functions
func errPrinter(err error) {
	if err != nil {
		color.Set(color.FgRed)
		fmt.Fprintf(os.Stderr, "[Error] %s\n", err.Error())
		color.Unset()
	}
}

func tcpStart(from string, to string) {
	proto := "tcp"

	localAddress, err := net.ResolveTCPAddr(proto, from)
	errHandler(err)

	remoteAddress, err := net.ResolveTCPAddr(proto, to)
	errHandler(err)

	listener, err := net.ListenTCP(proto, localAddress)
	errHandler(err)

	defer listener.Close()

	color.Set(color.FgGreen)
	fmt.Printf("Forwarding %s traffic from '%v' to '%v'\n", proto, localAddress, remoteAddress)
	color.Unset()

	for {
		src, err := listener.Accept()
		errHandler(err)
		fmt.Printf("New connection established from '%v'\n", src.RemoteAddr())
		go fwd(src, to, proto)
	}
}

func udpStart(from string, to string) {
	proto := "udp"

	localAddress, err := net.ResolveUDPAddr(proto, from)
	errHandler(err)

	remoteAddress, err := net.ResolveUDPAddr(proto, to)
	errHandler(err)

	listener, err := net.ListenUDP(proto, localAddress)
	errHandler(err)
	defer listener.Close()

	dst, err := net.DialUDP(proto, nil, remoteAddress)
	errHandler(err)
	defer dst.Close()

	color.Set(color.FgGreen)
	fmt.Printf("Forwarding %s traffic from '%v' to '%v'\n", proto, localAddress, remoteAddress)
	color.Unset()

	buf := make([]byte, 512)
	for {
		rnum, err := listener.Read(buf[0:])
		errHandler(err)

		_, err = dst.Write(buf[:rnum])
		errHandler(err)

		fmt.Printf("%d bytes forwared\n", rnum)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "fwd"
	app.Version = "1.0.1"
	app.Usage = "The little forwarder that could"
	app.UsageText = "fwd --from localhost:2222 --to 192.168.1.254:22"
	app.Copyright = "MIT License"
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Joel Bastos",
			Email: "kintoandar@gmail.com",
		},
	}
	app.Flags = []cli.Flag{
		cli.StringSliceFlag{
			Name:   "from, f",
			Value:  &cli.StringSlice{},
			EnvVar: "FWD_FROM",
			Usage:  "source HOST:PORT",
		},
		cli.StringSliceFlag{
			Name:   "to, t",
			EnvVar: "FWD_TO",
			Usage:  "destination HOST:PORT",
		},
		cli.BoolFlag{
			Name:  "list, l",
			Usage: "list local addresses",
		},
		cli.BoolFlag{
			Name:  "udp, u",
			Usage: "enable udp forwarding (tcp by default)",
		},
		cli.BoolFlag{
			Name:  "build, b",
			Usage: "build information",
		},
		cli.StringFlag{
			Name:  "config, c",
			Usage: "tunnels config file (YAML format), eg. ./tunnels.yaml",
		},
	}
	app.Action = func(c *cli.Context) error {
		defer color.Unset()
		color.Set(color.FgGreen)
		if c.Bool("list") {
			list, err := getLocalAddrs()
			errHandler(err)
			fmt.Println("Available local addresses:")
			color.Unset()
			for _, ip := range list {
				fmt.Println(ip)
			}
			return nil
		} else if c.Bool("build") {
			fmt.Println("Built with " + runtime.Version() + " for " + runtime.GOOS + "/" + runtime.GOARCH)
			color.Unset()
			return nil

		} else if c.String("to") == "" {
			color.Unset()
			cli.ShowAppHelp(c)
			return nil
		} else {
			tunnels := make([]*Tunnel, 0)

			configFile := c.String("config")
			if configFile != "" {
				config, err := readConfig(configFile)
				if err != nil {
					// just emit warning
					log.Printf("WARNING: %s", err.Error())
				}

				if config != nil {
					tunnels = config.Tunnels
				}
			}

			fromSl := c.StringSlice("from")
			toSl := c.StringSlice("to")

			if len(fromSl) < len(toSl) {
				return fmt.Errorf("invalid forwarding rules, [from] addresses are less than [to] addresses")
			}

			if len(fromSl) > len(toSl) {
				// in that case, we will pad toSl with the last one in order to match with fromSl length
				for {
					toSl = append(toSl, toSl[len(toSl)-1])
					if len(fromSl) == len(toSl) {
						break
					}
				}
			}

			for i := 0; i < len(fromSl); i++ {
				tunnel := Tunnel{
					Source: fromSl[i],
					Addr: toSl[i],
				}

				if c.Bool("udp") {
					tunnel.Protocol = protocolUDP
				} else {
					tunnel.Protocol = protocolTCP
				}

				tunnels = append(tunnels, &tunnel)
			}

			control := NewController(context.Background(), tunnels)

			return control.Run()
		}
	}
	app.Run(os.Args)
}
