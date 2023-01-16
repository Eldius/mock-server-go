package server

import (
    "fmt"
    "github.com/Eldius/mock-server-go/internal/mapper"
    "net/http"
)

func StartMockServer(port int, r *mapper.Router) {
    mux := http.NewServeMux()
    mux.HandleFunc("/", r.Handle)

    log.Infof("Starting mock server on port %d\n", port)
    log.WithError(http.ListenAndServe(fmt.Sprintf(":%d", port), mux)).Error("Failed to start HTTP server")
}
