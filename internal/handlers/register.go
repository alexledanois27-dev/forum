package handlers

import (
	"fmt"
	"net/http"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Page d'inscription")
}
