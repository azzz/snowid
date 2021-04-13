package config

import (
	"errors"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

type Config struct {
	machineID string
	epoch     string
	logLevel  string
	listen    string
}

func (c Config) MachineID() (uint64, error) {
	return strconv.ParseUint(c.machineID, 10, 64)
}

func (c Config) Listen() string {
	return c.listen
}

func (c Config) LogLevel() (logrus.Level, error) {
	return logrus.ParseLevel(c.logLevel)
}

func (c Config) Epoch() (time.Time, error) {
	return time.Parse("20060102150405", c.epoch)
}

// NewFromEnv creates a new Config from the environment variables
func NewFromEnv() (Config, error) {
	c := Config{}

	c.machineID = os.Getenv("MACHINE_ID")
	c.epoch = os.Getenv("EPOCH")
	c.logLevel = os.Getenv("LOG_LEVEL")
	c.listen = os.Getenv("LISTEN")

	if c.machineID == "" {
		return Config{}, errors.New("missing MACHINE_ID env")
	}

	c.logLevel = os.Getenv("LOG_LEVEL")
	if c.logLevel == "" {
		c.logLevel = "info"
	}

	return c, nil
}
