/*
 * Copyright (c) 2015, Chris Benedict <chrisbdaemon@gmail.com>
 * All rights reserved.
 *
 * Licensing terms are located in LICENSE file.
 */

package validate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInt(t *testing.T) {
	var err error

	err = Int("12")
	assert.NoError(t, err)

	err = Int("abc")
	assert.Error(t, err)
}

func TestPort(t *testing.T) {
	var err error

	// valid
	err = Port("8080")
	assert.NoError(t, err)

	// too low
	err = Port("-12")
	assert.Error(t, err)

	// too high
	err = Port("92162")
	assert.Error(t, err)

	// invalid int
	err = Port("two")
	assert.Error(t, err)
}

func TestHost(t *testing.T) {
	var err error

	// valid
	err = Host("8.8.8.8")
	assert.NoError(t, err)

	// first octet too high
	err = Host("256.21.64.126")
	assert.Error(t, err)
}
