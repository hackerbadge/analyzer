package main

import (
	"flag"
	"os"

	"github.com/dropbox/godropbox/errors"
	"github.com/yvasiyarov/globalconf"
)

type Config struct {
	Host         string
	Port         int
	RulesFile    string
	CollectorApi string
	Source       string
	DefaultXp    float64
}

func NewConfig(configPath string) (*Config, error) {
	c := &Config{}
	c.FlagVariables()
	err := c.LoadFile(configPath)
	return c, err
}

func (c *Config) FlagVariables() {
	flag.StringVar(&c.Host, "host", "localhost", "Web server host")
	flag.IntVar(&c.Port, "port", 3000, "port number")
	flag.StringVar(&c.RulesFile, "rules", "", "file with rules")
	flag.StringVar(&c.CollectorApi, "collector_api", "localhost:10100", "Collector API endpoint")
	flag.StringVar(&c.Source, "source", "github", "default source for promotions")
	flag.Float64Var(&c.DefaultXp, "default_xp", 2.0, "default XP for promotions")
}

func (c *Config) LoadFile(configPath string) error {
	options := &globalconf.Options{}

	if configPath != "" {
		if _, err := os.Stat(configPath); err != nil {
			return errors.New("can't load conf: " + err.Error())
		}
		options.Filename = configPath
	}

	// read config
	conf, err := globalconf.NewWithOptions(options)
	if err != nil {
		return errors.Wrap(err, "can't load conf")
	}
	conf.ParseAll()

	return nil
}
