package todo

import (
	"fmt"
	"html"
	"net/http"
	"time"
)

type todo struct {
	title      string
	status     string
	createdAt  time.Time
	finishedAt time.Time
}

type handler struct {
}

func New() *handler {
	return &handler{}
}
func (t *handler) AddTask(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Current time: ", time.Now())
	fmt.Fprintln(w, "URL Path: ", html.EscapeString(r.URL.Path))
	fmt.Fprintln(w, "HTTP Method: ", r.Method)
}
