package main

import (
	"net/http"

	"github.com/appointment-scheduling/cmd/handler"
)

func main() {
	h := handler.New()
	http.ListenAndServe(":8086", h.Router)
}
