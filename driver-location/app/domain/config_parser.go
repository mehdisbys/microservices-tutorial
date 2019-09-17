package domain

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/go-playground/validator"
	"gopkg.in/yaml.v2"
)

type Config struct {
	QueuePort    int    `yaml:"queue-port" validate:"required"`
	QueueHost    string `yaml:"queue-host" validate:"required"`
	QueueTopic   string `yaml:"queue-topic" validate:"required"`
	DatabasePort int    `yaml:"database-port" validate:"required"`
	DatabaseHost string `yaml:"database-host" validate:"required"`
}

// NewConfig returns a new `*Config` or an error if config file has missing and required values
func NewConfig(filename string) (*Config, error) {
	v := validator.New()

	cfg, err := GetServiceConfig(v, filename)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// GetServiceConfig get the service config as specified per struct at top of file
func GetServiceConfig(v *validator.Validate, filename string) (*Config, error) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	c := Config{}

	err = yaml.Unmarshal(source, &c)
	if err != nil {
		return nil, err
	}

	// Validate all required elements are in config
	err = v.Struct(c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
