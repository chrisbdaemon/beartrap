/*
 * Copyright (c) 2015, Chris Benedict <chrisbdaemon@gmail.com>
 * All rights reserved.
 *
 * Licensing terms are located in LICENSE file.
 */

package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// Params holds a list of all the parameters for an element in the configuration file
type Params map[string]string

// Config contains the datails of the configuration file
type Config struct {
	filename string
	data     []byte
	file     *os.File
}

// New takes in a filename, reads the contents into a Config object which is then returned
func New(filename string) (*Config, error) {
	cfg := new(Config)

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	cfg.data = data
	cfg.filename = filename

	return cfg, nil
}

// TrapParams returns an array with Params objects populated from the config file for traps
func (cfg *Config) TrapParams() ([]Params, error) {
	paramStruct := struct {
		Traps []Params
	}{nil}

	err := yaml.Unmarshal(cfg.data, &paramStruct)
	if err != nil {
		return nil, fmt.Errorf("Unable to process traps: %s", err)
	}

	return paramStruct.Traps, nil
}

// HandlerParams returns an array with Params objects populated from the config file for handlers
func (cfg *Config) HandlerParams() ([]Params, error) {
	paramStruct := struct {
		Handlers []Params
	}{nil}

	err := yaml.Unmarshal(cfg.data, &paramStruct)
	if err != nil {
		return nil, fmt.Errorf("Unable to process handlers: %s", err)
	}

	return paramStruct.Handlers, nil
}
