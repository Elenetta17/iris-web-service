package httpapi

import (
	"fmt"
	"net/http"
)

func FormPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `
        <form action="/hello" method="POST">
            <input name="name" placeholder="Your name">
            <button type="submit">Say Hello</button>
        </form>
    `)
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}
	name := r.PostFormValue("name")
	if name == "" {
		name = "World"
	}
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Hello %s!", name)
}
