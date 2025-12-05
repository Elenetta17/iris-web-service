package httpapi

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"strings"
)

//go:embed templates/*.html
var templateFS embed.FS

var templates = template.Must(template.ParseFS(templateFS, "templates/*.html"))

type HelloData struct {
	Name string
}

func FormPage(w http.ResponseWriter, r *http.Request) {
	log.Printf("FormPage called: %s %s", r.Method, r.URL.Path)
	templates.ExecuteTemplate(w, "form.html", nil)
	// No error check - if this fails, templates are broken (caught at startup)
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("HelloHandler called: %s %s", r.Method, r.URL.Path)

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	name := strings.TrimSpace(r.PostFormValue("name"))
	if name == "" {
		name = "World"
	}

	data := HelloData{
		Name: name,
	}

	templates.ExecuteTemplate(w, "hello.html", data)
}
