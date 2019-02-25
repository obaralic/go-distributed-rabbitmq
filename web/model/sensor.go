// -----------------------------------------------------------------------------
// Model package used for encapsulating web model.
// -----------------------------------------------------------------------------
package model

// -----------------------------------------------------------------------------
// Sensor - struct that encapsulates json model.
// -----------------------------------------------------------------------------
type Sensor struct {
	Name         string  `json:"name"`
	SerialNo     string  `json:"serialNo"`
	UnitType     string  `json:"unitType"`
	MinSafeValue float64 `json:"minSafeValue"`
	MaxSafeValue float64 `json:"maxSafeValue"`
}

// -----------------------------------------------------------------------------
// GetSensorByName - Gets sensor by name if exists.
// -----------------------------------------------------------------------------
func GetSensorByName(name string) (Sensor, error) {
	result := Sensor{}
	row := database.QueryRow(SELECT_SENSOR_ROW, name)
	err := row.Scan(
		&result.Name,
		&result.SerialNo,
		&result.UnitType,
		&result.MinSafeValue,
		&result.MaxSafeValue)

	return result, err
}
