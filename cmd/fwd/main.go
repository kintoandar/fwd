package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"

	"net/http"
	_ "net/http/pprof"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kintoandar/fwd/pkg/nethelp"
	"github.com/kintoandar/fwd/pkg/version"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
)

const (
	logFormatLogfmt = "logfmt"
	logFormatJSON   = "json"
)

type fwdConfig struct {
	LogLevel  string `yaml:"log.level"`
	LogFormat string `yaml:"log.format"`
	Fwd       []struct {
		From string `yaml:"from"`
		To   string `yaml:"to"`
	} `yaml:"fwd"`
}

func main() {
	// Command line interface configuration
	app := kingpin.New(filepath.Base(os.Args[0]), "fwd - The little forwarder that could")
	app.Help = `Name:
	fwd - 🚂 The little forwarder that could

Author:
  Joel Bastos @kintoandar

Example:
  fwd --from localhost:8000 --to example.com:80

`
	app.Author("Joel Bastos @kintoandar")

	app.HelpFlag.Short('h')

	from := app.Flag("from", "Local address to bind port (host:port)").
		Default("127.0.0.1:8000").
		OverrideDefaultFromEnvar("FWD_FROM").
		Short('f').
		String()

	to := app.Flag("to", "Remote address to forward traffic (host:port)").
		OverrideDefaultFromEnvar("FWD_TO").
		Short('t').String()

	appVersion := app.Flag("version", "Version details").
		Short('v').
		Bool()

	listNet := app.Flag("list", "List local network addresses available").
		Short('l').
		Bool()

	pprof := app.Flag("pprof", "Enable server runtime profiling data").Hidden().Bool()

	logLevel := app.Flag("log.level", "Logging level (error, warn, info, debug)").
		Default("info").
		Enum("error", "warn", "info", "debug")

	logFormat := app.Flag("log.format", "Logging format (logfmt, json)").
		Default(logFormatLogfmt).
		Enum(logFormatLogfmt, logFormatJSON)

	configFile := app.Flag("config", "Configuration file path (overrides all flags)").ExistingFile()

	_, err := app.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error parsing command line arguments")
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(2)
	}

	var config fwdConfig

	if *configFile != "" {
		data, err := ioutil.ReadFile(*configFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading config file")
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(2)
		}

		err = yaml.Unmarshal(data, &config)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error parsing config file")
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(2)
		}

		fmt.Println("log.format:", config.LogFormat)
		fmt.Println("log.level", config.LogLevel)
		*logLevel = config.LogLevel
		*logFormat = config.LogFormat
	}

	var logger log.Logger
	{
		var verbosity level.Option
		switch *logLevel {
		case "error":
			verbosity = level.AllowError()
		case "warn":
			verbosity = level.AllowWarn()
		case "info":
			verbosity = level.AllowInfo()
		case "debug":
			verbosity = level.AllowDebug()
		default:
			panic("undefined logging level")
		}
		logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
		if *logFormat == logFormatJSON {
			logger = log.NewJSONLogger(log.NewSyncWriter(os.Stderr))
		}
		logger = level.NewFilter(logger, verbosity)

		logger = log.With(
			logger,
			"ts",
			log.DefaultTimestampUTC,
			"caller",
			log.DefaultCaller)
	}

	if *pprof {
		level.Info(logger).Log(
			"msg", "starting pprof",
			"address", "https://localhost:6660/debug/pprof")

		go http.ListenAndServe("localhost:6660", nil)
	}

	switch {
	case *appVersion:
		fmt.Println("Version:", version.Version)
		fmt.Println("Go Version:", version.GoVersion)
		fmt.Println("OS / Arch:", version.OsArch)
		os.Exit(0)
	case *listNet:
		localAddr, err := nethelp.GetLocalAddrs()
		if err != nil {
			level.Error(logger).Log("msg", err.Error())
			os.Exit(1)
		}

		for _, ip := range localAddr {
			fmt.Println(ip)
		}

		os.Exit(0)

	case *configFile != "":
		var wg sync.WaitGroup

		for _, item := range config.Fwd {
			level.Info(logger).Log(
				"msg", "starting new fwd thread",
				"source", item.From,
				"destination", item.To)
			wg.Add(1)

			go runConn(logger, &wg, item.From, item.To)
		}

		wg.Wait()
	case *to != "":
		var wg sync.WaitGroup

		level.Info(logger).Log(
			"msg", "starting new fwd thread",
			"source", *from,
			"destination", *to)
		wg.Add(1)

		go runConn(logger, &wg, *from, *to)
		wg.Wait()
	default:
		app.Usage(os.Args[1:])
		os.Exit(2)
	}
}

func interrupt(logger log.Logger) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		level.Info(logger).Log(
			"msg", "signal caught, execution stopping",
			"signal", sig)
		os.Exit(0)
	}()
}

func fwd(logger log.Logger, dst io.Writer, src io.Reader) {
	_, err := io.Copy(dst, src)
	if err != nil {
		level.Error(logger).Log("msg", err.Error())
	}
}

func stream(logger log.Logger, src io.ReadWriter, remote string, proto string) {
	dst, err := net.Dial(proto, remote)
	if err != nil {
		level.Error(logger).Log("msg", err.Error())
	}

	go fwd(logger, src, dst)
	go fwd(logger, dst, src)
}

func runConn(logger log.Logger, wg *sync.WaitGroup, from string, to string) {
	defer wg.Done()

	localAddress, err := net.ResolveTCPAddr("tcp", from)
	if err != nil {
		level.Error(logger).Log("msg", err.Error())
		os.Exit(1)
	}

	remoteAddress, err := net.ResolveTCPAddr("tcp", to)
	if err != nil {
		level.Error(logger).Log("msg", err.Error())
		os.Exit(1)
	}

	listener, err := net.ListenTCP("tcp", localAddress)
	if err != nil {
		level.Error(logger).Log("msg", err.Error())
		os.Exit(1)
	}

	level.Info(logger).Log(
		"msg", "starting fwd",
		"from", localAddress,
		"to", remoteAddress,
		"protocol", "tcp")

	interrupt(logger)

	for {
		src, err := listener.Accept()
		if err != nil {
			level.Error(logger).Log("msg", err.Error())
			os.Exit(1)
		}

		level.Info(logger).Log(
			"msg", "new connection established",
			"source", src.RemoteAddr(),
			"destination", remoteAddress)

		stream(logger, src, to, "tcp")
	}
}
