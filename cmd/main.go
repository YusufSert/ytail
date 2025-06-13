package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path"
	"time"
	"ytail/client"
	"ytail/config"
	"ytail/tailer"
)

// todo: leart how config commnad work on other docker imagesj
// todo: ben parça parça module yapop bir leştiriyorum ana bir app olması gerekmez mi tail ve client birşeitirği ???
func main() {
	confPath := flag.String("config_path", "", "absolute path to configuration file")
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	/*
		wd, err := os.Getwd()
		if err != nil {
			logger.Error("couldn't get working directory", "err", err)
			return
		}
	*/
	globalConfig, err := config.ParseConfig(path.Join(*confPath, "ytail-config.yaml"))
	if err != nil {
		logger.Error("couldn't parse global config", "err", err)
		return
	}

	c := client.NewWithOptions(globalConfig.ClientConfig, client.WithLogger(logger))

	t, err := tailer.NewWithOptions(globalConfig.TailerConfig, tailer.WithClient(c), tailer.WithLogger(logger))
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		file, err := os.OpenFile(globalConfig.TailerConfig.ScrapePath+"/log.txt", os.O_WRONLY, os.ModeAppend)
		if err != nil {
			logger.Error(err.Error(), "", file)
			return
		}
		for {
			time.Sleep(1 * time.Second)

			_, err = fmt.Fprintf(file, `{"time":"2025-04-03T15:55:09.343298557+03:00","level":"ERROR","msg":"log-tail: error writing record to loki","error":"Post \"http://localhost:3100/loki/api/v1/push\": dial tcp [::1]:3100: connect: connection refused"}`)
			if err != nil {
				logger.Error(err.Error(), "", file)
			}
		}

	}()

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(10 * time.Second)
		cancel()
	}()
	err = t.Run(ctx)
	fmt.Println(err)
}
