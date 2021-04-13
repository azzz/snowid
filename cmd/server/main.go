package main

import (
	"github.com/azzz/ratatoskr/internal/config"
	"github.com/azzz/ratatoskr/internal/httpsrv"
	"github.com/azzz/ratatoskr/internal/sequence"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

func main() {
	var (
		logger    = logrus.New()
		timer     sequence.Timer
		machineID uint64
	)

	conf, err := config.NewFromEnv()
	if err != nil {
		logger.Errorf("configuration error: %s\n", err)
		os.Exit(1)
	}

	if lvl, err := conf.LogLevel(); err != nil {
		logger.Panicf("failed to parse log level: %s", err)
		os.Exit(1)
	} else {
		logger.Level = lvl
	}

	if epoch, err := conf.Epoch(); err != nil {
		logger.Panicf("failed to parse epoch: %s", err)
		os.Exit(1)
	} else {
		timer = func() uint64 {
			return uint64(time.Now().Sub(epoch) / 1_000_000)
		}
	}

	if m, err := conf.MachineID(); err != nil {
		logger.Panicf("failed to parse machine id: %s", err)
		os.Exit(1)
	} else {
		machineID = m
	}

	seq, err := sequence.NewSeq64(machineID, timer)
	if err != nil {
		logger.Panicf("Failed to initialize sequence: %s", err)
		os.Exit(1)
	}

	httpServer := httpsrv.New(logger.WithField("service", "http"), seq, conf.Listen())

	logger.Infof("Machine ID = %d", machineID)
	logger.Infof("Current time: %d", timer())
	logger.Infof("Starting http server on %s", conf.Listen())

	done := make(chan struct{})

	// run server in goroutine to let running another interface if a need be
	go func() {
		if err := httpServer.Start(); err != nil {
			logger.Panic("failed to start HTTP server: %s", err)
			os.Exit(1)
		}
	}()

	<-done
}
