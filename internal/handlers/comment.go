package handlers

import (
	"fmt"
	"net/http"
)

func CommentHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Page de commentaire")
}
