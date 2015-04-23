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
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/chrisbdaemon/beartrap/config"
	"github.com/chrisbdaemon/beartrap/config/validate"
)

// Trap holds data common to all trap types
type Trap struct {
	Severity int
}

// New takes in a params object and returns a trap
func New(params config.Params) (*Trap, error) {
	trap := new(Trap)

	errors := validParams(params)
	if len(errors) > 0 {
		for _, err := range errors {
			log.Println(err)
		}
		return nil, fmt.Errorf("Invalid trap parameters")
	}

	trap.Severity, _ = strconv.Atoi(params["severity"])

	return trap, nil
}

// validParams validates the paramters common to all traps.
func validParams(params config.Params) []error {
	errors := []error{}

	typeStr := strings.TrimSpace(params["type"])
	if len(typeStr) < 1 {
		errors = append(errors, fmt.Errorf("Missing trap type"))
	}

	err := validate.Int(params["severity"])
	if err != nil {
		errors = append(errors, fmt.Errorf("Invalid severity: %s", err))
	}

	return errors
}
