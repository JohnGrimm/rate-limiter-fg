package handler

import (
	"fmt"
	"net/http"
)

func HandlerHello(w http.ResponseWriter, rq *http.Request) {
	fmt.Fprintf(w, "Hello")
}
