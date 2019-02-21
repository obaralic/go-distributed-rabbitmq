// -----------------------------------------------------------------------------
// Storage package used for encapsulating work with persistence modueles.
// -----------------------------------------------------------------------------
package storage

import (
	"database/sql"
	"godistributed-rabbitmq/common"
	"log"
	"strconv"

	_ "github.com/lib/pq" // To initialize database drivers.
)

const (
	DATABASE_DRIVER = "postgres"
	DATABASE_SOURCE = "postgres://goadmin:goadmin@localhost/gosensors?sslmode=disable"
	DATABASE_CONSTR = "user=goadmin password=goadmin dbname=gosensors sslmode=disable"
)

// SQL Constants
const (
	CREATE_SENSORS_TABLE = `
		CREATE TABLE IF NOT EXISTS sensors (
			id SERIAL PRIMARY KEY,
      name varchar NOT NULL,
      serial_no varchar NOT NULL,
      unit_type varchar NOT NULL,
      max_safe_value float8 NOT NULL,
      min_safe_value float8 NOT NULL)
	`

	CREATE_READOUTS_TABLE = `
		CREATE TABLE IF NOT EXISTS readouts (
			id SERIAL PRIMARY KEY,
			value float8 NOT NULL,
			sensor_id integer,
			taken_on timestamp with time zone)
	`

	CLEAN_SENSORS_TABLE = `
		DELETE FROM sensors
	`

	INSERT_SENSOR = `
		INSERT INTO sensors (name, serial_no, unit_type, max_safe_value, min_safe_value)
			VALUES ($1, $2, $3, $4, $5)
	`

	INSERT_READOUT = `
		INSERT INTO readouts (value, sensor_id, taken_on)
			VALUES ($1, $2, $3)
	`
	SELECT_SENSORS = `
		SELECT id, name FROM sensors
	`
)

var sensorInit = [...][5]string{
	{"boiler_pressure_out", "MPR-728", "MPa", "15.4", "15.1"},
	{"condensor_pressure_out", "MPR-317", "MPa", "0.0022000000000000001", "0.00080000000000000004"},
	{"turbine_pressure_out", "MPR-492", "MPa", "1.3999999999999999", "0.80000000000000004"},
	{"boiler_temp_out", "XTLR-145", "C", "625", "580"},
	{"turbine_temp_out", "XTLR-145", "C", "115", "98"},
	{"condensor_temp_out", "XTLR-145", "C", "98", "83"},
}

// Database connection
var database *sql.DB

func init() {
	log.Println("Initializing PostgreSQL connector module...")
	var error error
	database, error = sql.Open(DATABASE_DRIVER, DATABASE_CONSTR)
	common.FailOnError(error, "Cannot open database")
	createTables()
	insertSensors()
}

// -----------------------------------------------------------------------------
// createTables - Creates database tables.
// -----------------------------------------------------------------------------
func createTables() {

	_, error := database.Exec(CREATE_SENSORS_TABLE)
	common.FailOnError(error, "Cannot create sensors table")

	_, error = database.Exec(CREATE_READOUTS_TABLE)
	common.FailOnError(error, "Cannot create readouts table")
}

// -----------------------------------------------------------------------------
// insertSensors - Cleans current sensors table and repopulates it.
// -----------------------------------------------------------------------------
func insertSensors() {
	_, error := database.Exec(CLEAN_SENSORS_TABLE)
	common.FailOnError(error, "Cannot clean sensors")

	for _, sensor := range sensorInit {
		maxSafeValue, _ := strconv.ParseFloat(sensor[3], 32)
		minSafeValue, _ := strconv.ParseFloat(sensor[4], 32)

		_, error = database.Exec(
			INSERT_SENSOR, sensor[0], sensor[1], sensor[2], maxSafeValue, minSafeValue)
		common.FailOnError(error, "Cannot insert sensor "+sensor[0])
	}
}

// -----------------------------------------------------------------------------
// GetSensors - Gets the map of sensors.
// -----------------------------------------------------------------------------
func GetSensors() (sensors map[string]int) {
	rows, error := database.Query(SELECT_SENSORS)
	common.FailOnError(error, "Cannot query sensors")

	sensors = make(map[string]int)
	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		sensors[name] = id
	}

	return
}
