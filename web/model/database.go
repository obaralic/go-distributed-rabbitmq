// -----------------------------------------------------------------------------
// Model package used for encapsulating web model.
// -----------------------------------------------------------------------------
package model

import (
	"database/sql"
	"godistributed-rabbitmq/common"
	"godistributed-rabbitmq/storage"

	_ "github.com/lib/pq"
)

const (
	SELECT_SENSOR_ROW = `
		SELECT name, serial_no, unit_type, min_safe_value, max_safe_value
			FROM sensors
				WHERE name = $1
	`
)

var database *sql.DB

func init() {
	var error error
	database, error = sql.Open(storage.DATABASE_DRIVER, storage.DATABASE_CONSTR)
	common.FailOnError(error, "Cannot open database")
}
