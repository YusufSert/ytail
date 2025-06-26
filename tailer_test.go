package ytail

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"testing"
	"time"
	"ytail/config"
	"ytail/pkg/client"
	tailer2 "ytail/pkg/tailer"
)

func TestTailer(t *testing.T) {
	level := &slog.LevelVar{}
	level.Set(slog.LevelDebug)
	opts := &slog.HandlerOptions{
		Level: level,
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))
	/*workingDir, err := os.Getwd()
	if err != nil {
		logger.Error("couldn't get working directory", "err", err)
		return
	}
	*/
	globalConfig, err := config.ParseFromFile("/home/yusufcan/Downloads/ytail/ytail-config.yaml")
	if err != nil {
		logger.Error("couldn't parse global config", "err", err)
		return
	}
	c := client.NewWithOptions(globalConfig.ClientConfig, client.WithLogger(logger))

	tr, err := tailer2.NewWithOptions(globalConfig.TailerConfig, tailer2.WithClient(c), tailer2.WithLogger(logger))
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		file, err := os.OpenFile(globalConfig.TailerConfig.ScrapePath+"/log-2025-06-19.txt", os.O_WRONLY, os.ModeAppend)
		if err != nil {
			logger.Error(err.Error(), "", file)
			return
		}
		for {
			time.Sleep(100 * time.Millisecond)

			_, err = fmt.Fprintf(file, `{"time":"2025-04-03T15:55:09.343298557+03:00","level":"ERROR","msg":"log-tail: error writing record to loki","error":"Post \"http://localhost:3100/loki/api/v1/push\": dial tcp [::1]:3100: connect: connection refused"}`+"\n")
			if err != nil {
				logger.Error(err.Error(), "", file)
			}
		}

	}()

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(100 * time.Second)
		cancel()
	}()
	err = tr.Run(ctx)
	t.Fatal(err)
}
