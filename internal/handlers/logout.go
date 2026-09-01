package handlers

import (
	"fmt"
	"net/http"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Page de déconnexion")
}
