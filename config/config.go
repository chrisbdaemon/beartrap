/*
 * Copyright (c) 2015, Chris Benedict <chrisbdaemon@gmail.com>
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * Redistributions of source code must retain the above copyright notice, this
 * list of conditions and the following disclaimer.
 *
 * Redistributions in binary form must reproduce the above copyright notice,
 * this list of conditions and the following disclaimer in the documentation
 * and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
 * LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
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
