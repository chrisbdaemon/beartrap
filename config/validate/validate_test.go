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

package validate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateInt(t *testing.T) {
	var err error

	err = ValidateInt("12")
	assert.NoError(t, err)

	err = ValidateInt("abc")
	assert.Error(t, err)
}

func TestValidatePort(t *testing.T) {
	var err error

	// valid
	err = ValidatePort("8080")
	assert.NoError(t, err)

	// too low
	err = ValidatePort("-12")
	assert.Error(t, err)

	// too high
	err = ValidatePort("92162")
	assert.Error(t, err)

	// invalid int
	err = ValidatePort("two")
	assert.Error(t, err)
}

func TestValidateHost(t *testing.T) {
	var err error

	// valid
	err = ValidateHost("8.8.8.8")
	assert.NoError(t, err)

	// first octet too high
	err = ValidateHost("256.21.64.126")
	assert.Error(t, err)
}
