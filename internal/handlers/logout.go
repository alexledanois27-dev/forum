package handlers

import "net/http"

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// TODO: middleware.DeleteSession(...)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
