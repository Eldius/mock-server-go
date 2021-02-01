package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/Eldius/mock-server-go/mapper"
)

var (
	tmpl         *template.Template
	adminConsole bool = true
)

func init() {
	if _, err := os.Stat("server/templates/index.html"); err == nil {
		tmpl = template.Must(template.ParseGlob("server/templates/*.html"))
	} else {
		fmt.Println("Admin console is disabled")
		adminConsole = false
	}
}

func AdminPanelHandler(router *mapper.Router) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		if !adminConsole {
			rw.WriteHeader(404)
			_, _ = rw.Write([]byte("Admin console is not available"))
			return
		}
		err := tmpl.ExecuteTemplate(rw, "index.html", router)
		if err != nil {
			log.Printf("Failed to parse template: %s\n", err)
		}
	}
}
