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
	"log"
	"net"
	"strconv"

	"github.com/chrisbdaemon/beartrap/config"
	"github.com/chrisbdaemon/beartrap/config/validate"
)

type trapInterface interface {
	Validate() []error
}

// Start opens a TCP socket and alerts on all connections
func (trap *TCPTrap) Start() error {
	l, err := net.Listen("tcp", ":"+strconv.Itoa(trap.port))
	if err != nil {
		log.Fatalf("Unable to open tcp port %d: %s", trap.port, err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && !ne.Temporary() {
				log.Println("Error accepting connection:", err)
			}
		}
		handleConnection(conn)
	}
}

// Use a slightly simplified connection interface
// to make testing a little easier
type simpleConn interface {
	Close() error
	RemoteAddr() net.Addr
}

func handleConnection(c simpleConn) {
	fmt.Println("Got connection:", c)
	c.Close()
}

// TCPTrap contains the details of a TCP trap
type TCPTrap struct {
	baseTrap trapInterface
	host     string
	port     int
	params   config.Params
}

// New returns a new trap object waiting to be engaged
func New(params config.Params, baseTrap trapInterface) *TCPTrap {
	tcptrap := new(TCPTrap)
	tcptrap.baseTrap = baseTrap

	tcptrap.params = params
	tcptrap.port, _ = strconv.Atoi(params["port"])
	tcptrap.host = params["host"]

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

	err = validate.Host(trap.params["host"])
	if err != nil {
		errors = append(errors, fmt.Errorf("Invalid host: %s", err))
	}

	return errors
}
