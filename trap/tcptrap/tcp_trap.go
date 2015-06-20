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

package tcptrap

import (
	"fmt"
	"strconv"

	"github.com/chrisbdaemon/beartrap/config"
	"github.com/chrisbdaemon/beartrap/config/validate"
)

type trapInterface interface {
	Validate() []error
}

// TCPTrap contains the details of a TCP trap
type TCPTrap struct {
	baseTrap trapInterface
	port     int
	params   config.Params
}

// New returns a new trap object waiting to be engaged
func New(params config.Params, baseTrap trapInterface) *TCPTrap {
	tcptrap := new(TCPTrap)
	tcptrap.baseTrap = baseTrap

	tcptrap.params = params
	tcptrap.port, _ = strconv.Atoi(params["port"])

	return tcptrap
}

// Validate performs validation on the trap and returns any errors
func (trap *TCPTrap) Validate() []error {
	var errors []error

	errors = trap.baseTrap.Validate()

	err := validate.Port(trap.params["port"])
	if err != nil {
		errors = append(errors, fmt.Errorf("Invalid port: %s", err))
	}

	return errors
}
