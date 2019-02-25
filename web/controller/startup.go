// -----------------------------------------------------------------------------
// Web package used for encapsulating web controllers.
// -----------------------------------------------------------------------------
package controller

import (
	"log"
	"net/http"
)

var sockets *SocketController

func init() {
	log.Printf("Startup initialize socket")
	sockets = NewSocketController()
}

// -----------------------------------------------------------------------------
// Initialize - On demand initializator function.
// -----------------------------------------------------------------------------
func Initialize() {
	registerRoutes()
	registerFileServers()
}

// -----------------------------------------------------------------------------
// registerRoutes - For handling web socket calls.
// -----------------------------------------------------------------------------
func registerRoutes() {
	http.HandleFunc("/ws", sockets.handleMessage)
}

// -----------------------------------------------------------------------------
// registerFileServers - For public assets and node modules.
// -----------------------------------------------------------------------------
func registerFileServers() {
	http.Handle("/public/", http.FileServer(http.Dir("assets")))

	http.Handle("/public/lib/",
		http.StripPrefix("/public/lib/", http.FileServer(http.Dir("node_modules"))))
}
