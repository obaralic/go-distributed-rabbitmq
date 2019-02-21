package storage

import (
	"errors"
	"godistributed-rabbitmq/common/dto"
)

// Maps sensor name to it's database id.
var sensors map[string]int

// -----------------------------------------------------------------------------
// SaveReadout - Saves the readout from the given sensor.
// -----------------------------------------------------------------------------
func SaveReadout(readout *dto.Readout) (error error) {
	if sensors[readout.Name] == 0 {
		sensors = GetSensors()
	}

	if sensors[readout.Name] == 0 {
		return errors.New("Unable to find sensor for name '" + readout.Name + "''")
	}

	_, error = database.Exec(INSERT_READOUT, readout.Value, sensors[readout.Name], readout.Timestamp)
	return
}
