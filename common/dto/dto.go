// -----------------------------------------------------------------------------
// This package represents Data Transfer Object used to transport sensor data
// from the producer on the one side to the consumer on the other.
// -----------------------------------------------------------------------------
package dto

import (
	"encoding/gob"
	"time"
)

// -----------------------------------------------------------------------------
// Readout - Represents sensor readout structure that is transfered further on.
// -----------------------------------------------------------------------------
type Readout struct {
	Name      string    // Name of the originating sensor.
	Value     float64   // Value of the sensor readout.
	Timestamp time.Time // Timestamp of the sensor readout.
}

func init() {
	// Register type with gob as a fast and efficient way of (de)serializing data.
	// Gob is recommended when both client and server are written in Go.
	// Also consider Cap’n Proto or Protocol Buffers.
	gob.Register(Readout{})
}

// -----------------------------------------------------------------------------
// NewReadout - Creates new readout for the goven values.
// -----------------------------------------------------------------------------
func NewReadout(name string, readout float64, timestamp time.Time) *Readout {
	return &Readout{
		Name:      name,
		Value:     readout,
		Timestamp: timestamp}
}
