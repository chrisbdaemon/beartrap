/*
 * Copyright (c) 2015, Chris Benedict <chrisbdaemon@gmail.com>
 * All rights reserved.
 *
 * Licensing terms are located in LICENSE file.
 */

package broadcast

import (
	"testing"

	"github.com/chrisbdaemon/beartrap/alert"
	"github.com/stretchr/testify/assert"
)

func TestAddReceiver(t *testing.T) {
	d := Broadcast{}
	c := make(chan alert.Alert)

	d.AddReceiver(c)
	assert.Equal(t, 1, len(d.receivers), "Unable to add receiver to Dispatch")

	// Ensure no duplicates
	d.AddReceiver(c)
	assert.Equal(t, 1, len(d.receivers), "AddReceiver should not permit duplicates")
}

func TestRemoveReceiver(t *testing.T) {
	d := Broadcast{}
	c := make(chan alert.Alert)

	d.receivers = []chan alert.Alert{c}
	d.RemoveReceiver(c)
	assert.Equal(t, 0, len(d.receivers), "Unable to remove receiver")
}

func TestBroadcastAlert(t *testing.T) {
	d := Broadcast{}
	c := make(chan alert.Alert, 1)

	d.receivers = []chan alert.Alert{c}
	a := alert.Alert{}

	d.BroadcastAlert(a)

	select {
	case _ = <-c:
	default:
		assert.Fail(t, "BroadcastAlert did not correctly send alert")
	}
}
