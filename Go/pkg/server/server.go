package server

import (
	"net/http"
)

func ServerStart() error {

	webDir := "./web"
	fileServer := http.FileServer(http.Dir(webDir))

	http.Handle("/", fileServer)

	return http.ListenAndServe(":7540", nil)
}
