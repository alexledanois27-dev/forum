package handlers

import "net/http"

func LikeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// TODO: Vérifier la session

	// TODO: Récupérer le post/commentaire concerné

	// TODO: Ajouter ou retirer un like

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
