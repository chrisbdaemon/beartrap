/*
 * Copyright (c) 2015, Chris Benedict <chrisbdaemon@gmail.com>
 * All rights reserved.
 *
 * Licensing terms are located in LICENSE file.
 */

package trap

import (
	"testing"

	"github.com/chrisbdaemon/beartrap/alert"
	"github.com/chrisbdaemon/beartrap/config"
	"github.com/stretchr/testify/assert"
)

type fakeDispatcher struct{}

func (d fakeDispatcher) BroadcastAlert(alert.Alert) {}

func TestNew(t *testing.T) {
	params := config.Params{
		"type":     "tcp",
		"severity": "3",
	}

	trap, _ := New(params, fakeDispatcher{})
	assert.NotNil(t, trap)
}

func TestValidate(t *testing.T) {
	trap := new(BaseTrap)
	trap.params = config.Params{
		"severity": "3",
	}

	errors := trap.Validate()
	assert.Equal(t, 0, len(errors))

	// return error for invalid severity
	trap.params["severity"] = "five"
	errors = trap.Validate()
	assert.Equal(t, 1, len(errors))
}
