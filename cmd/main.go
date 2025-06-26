package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"ytail/config"
	"ytail/pkg/client"
	"ytail/pkg/tailer"
)

// find way or check promtail to give precedence flag vars
func main() {
	confPath := flag.String("config.path", "", "absolute path to configuration file")
	logLevel := flag.Int("log.lvl", 0, "specifies the log level of the logger")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.Level(*logLevel),
		AddSource: true,
	}))

	globalConfig, err := config.ParseFromFile(*confPath)
	if err != nil {
		logger.Error("couldn't parse global config", "path", *confPath, "err", err)
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
