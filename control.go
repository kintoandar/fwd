package main

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
)

const protocolUDP = "udp"
const protocolTCP = "tcp"

type tunnelConfigFile struct {
	Tunnels []*Tunnel `yaml:"tunnels"`
}

func readConfig(filename string) (*tunnelConfigFile, error) {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to read tunnels config: %s", err)
	}

	var c tunnelConfigFile

	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %s", err)
	}

	return &c, nil
}

type Tunnel struct {
	Addr     string `json:"addr" yaml:"addr"`
	Protocol string `json:"protocol" yaml:"protocol"`
	Source   string `json:"source" yaml:"source"`
}

type Controller struct {
	ctx     context.Context
	Tunnels []*Tunnel
}

func NewController(ctx context.Context, tunnels []*Tunnel) *Controller {
	ctx, cancel := context.WithCancel(ctx)

	c := &Controller{
		ctx:     ctx,
		Tunnels: tunnels,
	}

	c.handleInterrupt(cancel)

	return c
}

func (c *Controller) printCtrlC() {
	color.Set(color.FgYellow)
	fmt.Println()
	fmt.Println("<CTRL+C> to exit")
	fmt.Println()
	color.Unset()
}

func (c *Controller) handleInterrupt(cancel context.CancelFunc) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		color.Set(color.FgGreen)
		fmt.Println("\nExecution stopped by", sig)
		color.Unset()
		cancel()
		os.Exit(0)
	}()
}

func (c *Controller) Run() error {
	c.printCtrlC()

	for _, tunnel := range c.Tunnels {
		switch tunnel.Protocol {
		case protocolTCP:
			go tcpStart(tunnel.Source, tunnel.Addr)
		case protocolUDP:
			go udpStart(tunnel.Source, tunnel.Addr)
		default:
			return fmt.Errorf("unsupported protocol '%s'", tunnel.Protocol)
		}
	}

	for {
		select {
		case <-c.ctx.Done():
			return fmt.Errorf("context ended, exiting")
		}
	}
}
