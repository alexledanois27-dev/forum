package handlers

import "net/http"

func PostHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:

		// TODO: Afficher un post

	case http.MethodPost:

		// TODO: Vérifier la session

		// TODO: Récupérer title/content

		// TODO: Validation

		// TODO: database.CreatePost(...)

		http.Redirect(w, r, "/", http.StatusSeeOther)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
