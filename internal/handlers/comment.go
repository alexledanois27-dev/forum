package handlers

import "net/http"

func CommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// TODO: Vérifier la session

	// TODO: Récupérer le commentaire

	// TODO: Validation

	// TODO: database.CreateComment(...)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
