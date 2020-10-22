package main

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"os"
	"os/signal"
	"syscall"
)

const protocolUDP = "upd"
const protocolTCP = "tcp"

type Tunnel struct {
	Addr     string `json:"addr"`
	Protocol string `json:"protocol"`
	Source   string `json:"source"`
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

	c.printCtrlC()

	for {
		select {
		case <-c.ctx.Done():
			return fmt.Errorf("context ended, exiting")
		}
	}
}
