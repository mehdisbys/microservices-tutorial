package domain

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/go-playground/validator"

	"gopkg.in/yaml.v2"
)

type Config struct {
	DriverHost              string        `yaml:"driver-service-host" validate:"required"`
	DriverLocationsEndpoint string        `yaml:"driver-locations-endpoint" validate:"required"`
	Timeout                 time.Duration `yaml:"timeout" validate:"required"`
	MinDistance             float64       `yaml:"minimum-distance" validate:"required"`
}

// NewConfig returns a new `*Config`
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

	err = v.Struct(c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
