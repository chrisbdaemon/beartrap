/*
 * Copyright (c) 2015, Chris Benedict <chrisbdaemon@gmail.com>
 * All rights reserved.
 *
 * Licensing terms are located in LICENSE file.
 */

package handler

import (
	"log"
	"testing"

	"github.com/chrisbdaemon/beartrap/alert"
	"github.com/chrisbdaemon/beartrap/config"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	c := make(chan alert.Alert)
	params := config.Params{
		"type":      "syslog",
		"threshold": "3",
	}

	handler, _ := New(params, c)
	assert.NotNil(t, handler)
}

func TestValidate(t *testing.T) {
	handler := new(BaseHandler)
	handler.params = config.Params{
		"threshold": "3",
	}

	errors := handler.Validate()
	log.Println(errors)
	assert.Equal(t, 0, len(errors))

	// return error for invalid severity
	handler.params["threshold"] = "five"
	errors = handler.Validate()
	assert.Equal(t, 1, len(errors))
}
