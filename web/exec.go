// -----------------------------------------------------------------------------
// Web package used for encapsulating work with web application.
// -----------------------------------------------------------------------------
package main

import (
	"godistributed-rabbitmq/web/controller"
	_ "godistributed-rabbitmq/web/controller"
	"net/http"
)

func main() {
	controller.Initialize()

	http.ListenAndServe(":3000", nil)
}
