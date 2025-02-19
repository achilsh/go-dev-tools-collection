package common

import (
	"log"
	"os"

	"go.uber.org/dig"
)

type Config struct {
	Prefix string
}

func ExampleCall() *dig.Container {

	c := dig.New()

	err := c.Provide(func() (*Config, error) {
		return &Config{Prefix: "[foo] "}, nil
	})
	if err != nil {
		panic(err)
	}
	err = c.Provide(func(cfg *Config) *log.Logger {
		return log.New(os.Stdout, cfg.Prefix, 0)
	})
	if err != nil {
		panic(err)
	}
	return c
}
