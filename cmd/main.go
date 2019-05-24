package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync_study/configuration"
	"sync_study/entity"
	"sync_study/logger"
	"sync_study/parse/study"
	"sync_study/scheduler"
	"syscall"
	"time"
)

func init() {
	flag.StringVar(&configuration.ConfigFile, "configfile", "./config.yaml", "config file path")
}

func main() {

	flag.Parse()

	configuration.Load()

	cfg := configuration.Config()
	err := logger.Init(cfg.Log.IsStdout, cfg.Log.LogFile, cfg.Log.Level)
	if err != nil {
		panic(fmt.Errorf("日志信息配置失败~%s", err))
	}

	logger.Info(configuration.ConfigFile)

	ctx, cancelFunc := context.WithCancel(context.Background())

	s := scheduler.NewScheduler(3, ctx)
	go s.Run()

	time.Sleep(3 * time.Second)

	s.Submit(
		entity.Request{
			URI:          cfg.DataURI["testUri"],
			ParseHandler: study.ParseStudy,
			RequestConfig: entity.RequestConfig{
				Retry:      5,
				RetryIndex: 0,
			},
			Data: &entity.StudyRequest{
				Depth: 0,
			},
		},
	)

	signals := make(chan os.Signal, 0)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	go func() {
		for {
			select {
			case <-signals:
				cancelFunc()
				logger.Info("ctrl+c out")
			}
		}
	}()

	select {
	case <-ctx.Done():
	}
}
