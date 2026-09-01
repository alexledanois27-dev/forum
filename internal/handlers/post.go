package handlers

import (
	"fmt"
	"net/http"
)

func PostHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Page post")
}
