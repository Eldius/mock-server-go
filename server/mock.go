package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Eldius/mock-server-go/mapper"
)

func StartMockServer(port int, r *mapper.Router) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", r.Handle)

	log.Println(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}
