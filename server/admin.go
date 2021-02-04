package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Eldius/mock-server-go/mapper"
	"gopkg.in/yaml.v3"
)

func RouteHandler(router *mapper.Router) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		log.Printf("Receiving %s request\n", r.Method)
		if r.Method == "POST" {
			var mapping mapper.RequestMapping
			err := json.NewDecoder(r.Body).Decode(&mapping)
			if err != nil {
				log.Println("Failed to read request body")
				log.Println(err.Error())
				rw.WriteHeader(500)
				_, _ = rw.Write([]byte(err.Error()))
				return
			}
			router.Add(mapping)
			rw.WriteHeader(200)
			rw.Header().Add("Content-Type", "application/json")
			_ = json.NewEncoder(rw).Encode(mapping)

		} else if r.Method == "GET" {
			rw.WriteHeader(200)
			encodeResponse(router, r, rw)
		} else {
			log.Println("returning: 'Method not allowed'")
			rw.WriteHeader(405)
		}
	}
}

func RequestsHandler(router *mapper.Router) func(rw http.ResponseWriter, r *http.Request) {

	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			log.Println("returning: 'Method not allowed'")
			rw.WriteHeader(405)
			return
		}

		rw.WriteHeader(200)
		rw.Header().Add("Content-Type", "application/json")
		_ = json.NewEncoder(rw).Encode(router.GetRequests())
	}
}

func encodeResponse(router *mapper.Router, r *http.Request, rw http.ResponseWriter) {
	accepts := r.Header.Get("Accept")
	if strings.Contains(strings.ToLower(accepts), "application/yaml") {
		rw.Header().Add("Content-Type", "application/yaml")
		_ = yaml.NewEncoder(rw).Encode(router)
	} else {
		rw.Header().Add("Content-Type", "application/json")
		_ = json.NewEncoder(rw).Encode(router)
	}
}

func StartAdminServer(port int, r *mapper.Router) {
	mux := http.NewServeMux()
	if adminConsole {
		fs := http.FileServer(http.Dir("./static"))
		mux.Handle("/static/", http.StripPrefix("/static/", fs))
		mux.HandleFunc("/", AdminPanelHandler(r))
	}
	mux.HandleFunc("/route", RouteHandler(r))
	mux.HandleFunc("/request", RequestsHandler(r))

	host := fmt.Sprintf(":%d", port)

	log.Println(http.ListenAndServe(host, mux))
}
