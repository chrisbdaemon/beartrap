/*
 * Copyright (c) 2015, Chris Benedict <chrisbdaemon@gmail.com>
 * All rights reserved.
 *
 * Licensing terms are located in LICENSE file.
 */

package broadcast

import (
	"fmt"

	"github.com/chrisbdaemon/beartrap/alert"
)

// Broadcast handles broadcasting alerts from traps to the handlers
type Broadcast struct {
	receivers []chan alert.Alert
}

// AddReceiver registers a new receiver to receive alerts. Ignores duplicates
func (b *Broadcast) AddReceiver(c chan alert.Alert) {
	// Prevent duplicates
	_ = "breakpoint"
	_, err := b.indexOfReceiver(c)
	if err == nil {
		return
	}

	b.receivers = append(b.receivers, c)
}

// RemoveReceiver removes a receiver
func (b *Broadcast) RemoveReceiver(c chan alert.Alert) {
	i, err := b.indexOfReceiver(c)
	if err != nil {
		return
	}

	b.receivers = append(b.receivers[:i], b.receivers[i+1:]...)
}

func (b *Broadcast) indexOfReceiver(c chan alert.Alert) (int, error) {
	for i := range b.receivers {
		if c == b.receivers[i] {
			return i, nil
		}
	}

	return -1, fmt.Errorf("Receiver not found")
}

// BroadcastAlert sends an alert object to all recievers added via AddReceiver()
// All reciever channels needs a listener or this will block until the receivers do
func (b Broadcast) BroadcastAlert(a alert.Alert) {
	for i := range b.receivers {
		b.receivers[i] <- a
	}
}
