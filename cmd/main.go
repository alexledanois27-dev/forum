package main

import (
	"forum/internal/database"
	"forum/internal/handlers"
	"log"
	"net/http"
)

func main() {
	db, err := database.Open("data/forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	mux := http.NewServeMux()

	h := &handlers.Handler{
		DB: db,
	}

	mux.HandleFunc("/", h.HomeHandler)
	mux.HandleFunc("/login", h.LoginHandler)
	mux.HandleFunc("/register", h.RegisterHandler)
	mux.HandleFunc("/logout", h.LogoutHandler)
	mux.HandleFunc("/post", h.PostHandler)
	mux.HandleFunc("/comment", h.CommentHandler)
	mux.HandleFunc("/like", h.LikeHandler)

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	log.Println("Server started on http://localhost:8080")

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
