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
