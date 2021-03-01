package server

import (
	"fmt"
	"net/http"

	"github.com/Eldius/mock-server-go/mapper"
)

func StartMockServer(port int, r *mapper.Router) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", r.Handle)

	log.Infof("Starting mock server on port %d\n", port)
	log.WithError(http.ListenAndServe(fmt.Sprintf(":%d", port), mux)).Error("Failed to start HTTP server")
}
