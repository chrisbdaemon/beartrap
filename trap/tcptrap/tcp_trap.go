/*
 * Copyright (c) 2015, Chris Benedict <chrisbdaemon@gmail.com>
 * All rights reserved.
 *
 * Licensing terms are located in LICENSE file.
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
