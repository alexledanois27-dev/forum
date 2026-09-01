package handlers

import (
	"fmt"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Page de connexion")
}
