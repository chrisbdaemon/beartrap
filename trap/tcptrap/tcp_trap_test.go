/*
 * Copyright (c) 2015, Chris Benedict <chrisbdaemon@gmail.com>
 * All rights reserved.
 *
 * Licensing terms are located in LICENSE file.
 */

package tcptrap

import (
	"net"
	"testing"

	"github.com/chrisbdaemon/beartrap/config"
	"github.com/stretchr/testify/assert"
)

var baseParams = config.Params{
	"type": "tcp",
	"port": "5555",
	"host": "127.0.0.1",
}

type stubBaseTrap struct{}

func (t stubBaseTrap) TriggerAlert(s string) {}
func (t stubBaseTrap) Validate() []error {
	var errors []error
	return errors
}

type stubConn struct{}

func (t stubConn) Close() error         { return nil }
func (t stubConn) RemoteAddr() net.Addr { a, _ := net.InterfaceAddrs(); return a[0] }

func TestNew(t *testing.T) {
	var baseTrap stubBaseTrap
	params := config.Params{}

	trap := New(params, baseTrap)
	assert.NotNil(t, trap)
}

func TestValidate(t *testing.T) {
	var baseTrap stubBaseTrap
	params := baseParams

	trap := New(params, baseTrap)
	errors := trap.Validate()
	assert.Equal(t, 0, len(errors))

	params["port"] = "-100"
	errors = trap.Validate()
	assert.Equal(t, 1, len(errors))

	params["host"] = "foobar"
	errors = trap.Validate()
	assert.Equal(t, 2, len(errors))
}

func TestHandleConnection(t *testing.T) {
	trap := stubTrap()
	c := stubConn{}

	// Just exercise the code for now
	trap.handleConnection(c)
}

func stubTrap() *TCPTrap {
	var baseTrap stubBaseTrap
	params := baseParams

	return New(params, baseTrap)
}
