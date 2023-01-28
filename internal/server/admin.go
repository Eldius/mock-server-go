package server

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/Eldius/mock-server-go/internal/logger"
	"github.com/Eldius/mock-server-go/internal/mapper"
	"io/fs"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

var (
	adminConsole = true
	log          = logger.Log()

	//go:embed static/**
	staticFiles embed.FS
)

func init() {

}

func RouteHandler(router *mapper.Router) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		log.Printf("Receiving %s request\n", r.Method)
		if r.Method == "POST" {
			var mapping mapper.RequestMapping
			err := json.NewDecoder(r.Body).Decode(&mapping)
			if err != nil {
				log.WithError(err).
					Error("Failed to read request body")
				rw.WriteHeader(500)
				_, _ = rw.Write([]byte(err.Error()))
				return
			}
			router.Add(mapping)
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(200)
			_ = json.NewEncoder(rw).Encode(mapping)

		} else if r.Method == "GET" {
			encodeResponse(router, r, rw)
		} else {
			log.Warn("returning: 'Method not allowed'")
			rw.WriteHeader(405)
		}
	}
}

func RequestsHandler(router *mapper.Router) func(rw http.ResponseWriter, r *http.Request) {

	return func(rw http.ResponseWriter, r *http.Request) {
		_log := log.WithFields(logrus.Fields{
			"path":   r.URL.Path,
			"method": r.Method,
		})
		if r.Method != "GET" {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			_log.WithFields(logrus.Fields{
				"code": http.StatusMethodNotAllowed,
			}).Warn("Method not allowed")
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		_log.WithFields(logrus.Fields{
			"code": http.StatusOK,
		}).Debug("OK")
		_ = json.NewEncoder(rw).Encode(router.GetRequests())
	}
}

func encodeResponse(router *mapper.Router, r *http.Request, rw http.ResponseWriter) {
	accepts := r.Header.Get("Accept")
	if strings.Contains(strings.ToLower(accepts), "application/yaml") {
		rw.Header().Set("Content-Type", "application/yaml")
		rw.WriteHeader(200)
		_ = yaml.NewEncoder(rw).Encode(router)
	} else {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(200)
		_ = json.NewEncoder(rw).Encode(router)
	}
}

func StartAdminServer(port int, r *mapper.Router) {
	mux := http.NewServeMux()
	if adminConsole {
		fsys := fs.FS(staticFiles)
		html, _ := fs.Sub(fsys, "static")
		fs := http.FileServer(http.FS(html))
		mux.Handle("/", fs)
	}
	mux.HandleFunc("/route", RouteHandler(r))
	mux.HandleFunc("/request", RequestsHandler(r))

	host := fmt.Sprintf(":%d", port)

	log.Infof("Starting admin server on port %d\n", port)
	log.WithError(http.ListenAndServe(host, mux)).Error("Failed to start HTTP server")
}
