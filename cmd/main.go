package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"ytail/config"
	"ytail/pkg/client"
	"ytail/pkg/tailer"
)

// find way or check promtail to give precedence flag vars
// add flag validation, and better description of flags and help flag --help, -h
func main() {
	confPath := flag.String("config.path", "", "absolute path to configuration file")
	logLevel := flag.Int("log.lvl", 0, "specifies the log level of the logger")
	help := flag.Bool("help", false, "Show this help message")
	flag.Parse()

	if *help {
		printUsage()
		os.Exit(0)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.Level(*logLevel),
		AddSource: true,
	}))

	globalConfig, err := config.ParseFromFile(*confPath)
	if err != nil {
		logger.Error("couldn't parse config", "path", *confPath, "err", err)
		os.Exit(1)
	}

	c := client.NewWithOptions(globalConfig.ClientConfig, client.WithLogger(logger))
	defer c.Stop()

	t, err := tailer.NewWithOptions(globalConfig.TailerConfig, tailer.WithClient(c), tailer.WithLogger(logger))
	if err != nil {
		logger.Error("couldn't start tailer", "err", err)
		os.Exit(1)
	}
	err = t.Run(context.Background())
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}

const version = "0.0.1"
const usage = `ytail is a tailer that read log lines from file and sends them to loki

Usage: ytail [options]

Version: %s

Options:
  -log.lvl,                	log level. available log levels:
                           	- -4 = DEBUG
                           	-  0 = INFO
                           	-  4	= WARN
                           	-  8 = ERROR
  -config.path,       		config file path.
`

func printUsage() {
	fmt.Printf(usage, version)
}
