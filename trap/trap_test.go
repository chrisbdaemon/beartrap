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

package trap

import (
	"testing"

	"github.com/chrisbdaemon/beartrap/config"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	params := config.Params{
		"type":     "tcp",
		"severity": "3",
	}

	trap, err := New(params)
	assert.NotNil(t, trap)
	assert.NoError(t, err)

	params["type"] = ""
	trap, err = New(params)
	assert.Nil(t, trap)
	assert.Error(t, err)
}

func TestValidate(t *testing.T) {
	params := config.Params{
		"type":     "tcp",
		"severity": "3",
	}

	errors := validParams(params)
	assert.Equal(t, 0, len(errors))

	// return error for invalid severity
	params["severity"] = "five"
	errors = validParams(params)
	assert.Equal(t, 1, len(errors))

	// return error for severity and type
	errors = validParams(config.Params{})
	assert.Equal(t, 2, len(errors))
}
