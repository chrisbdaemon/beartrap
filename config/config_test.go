/*
 * Copyright (c) 2015, Chris Benedict <chrisbdaemon@gmail.com>
 * All rights reserved.
 *
 * Licensing terms are located in LICENSE file.
 */

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	// Ensure it works with valid input
	filename := "test_conf_files/config_test_empty.yml"
	cfg, err := New(filename)
	assert.NotNil(t, cfg, "should create new object")
	assert.NoError(t, err, "should not return error")

	// Ensure it throws errors when needed
	filename = "test_conf_files/config_test_fake.yml"
	cfg, err = New(filename)
	assert.Nil(t, cfg, "should return nil")
	assert.Error(t, err, "should return error")
}

func TestTrapParams(t *testing.T) {
	// Valid config file
	cfg, _ := New("test_conf_files/config_test_traps.yml")
	trapParams, err := cfg.TrapParams()
	assert.Equal(t, 2, len(trapParams))
	assert.NoError(t, err)

	// Invalid config file
	cfg, _ = New("test_conf_files/config_test_traps_bad.yml")
	trapParams, err = cfg.TrapParams()
	assert.Equal(t, 0, len(trapParams))
	assert.Error(t, err)
}

func TestHandlerParams(t *testing.T) {
	// Valid config file
	cfg, _ := New("test_conf_files/config_test_handlers.yml")
	handlerParams, err := cfg.HandlerParams()
	assert.Equal(t, 2, len(handlerParams))
	assert.NoError(t, err)

	// Invalid config file
	cfg, _ = New("test_conf_files/config_test_handlers_bad.yml")
	handlerParams, err = cfg.HandlerParams()
	assert.Equal(t, 0, len(handlerParams))
	assert.Error(t, err)
}
