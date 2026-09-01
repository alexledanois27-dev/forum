package main

import (
	"forum/internal/handlers"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.HomeHandler)
	mux.HandleFunc("/login", handlers.LoginHandler)
	mux.HandleFunc("/register", handlers.RegisterHandler)
	mux.HandleFunc("/logout", handlers.LogoutHandler)
	mux.HandleFunc("/post", handlers.PostHandler)
	mux.HandleFunc("/comment", handlers.CommentHandler)
	mux.HandleFunc("/like", handlers.LikeHandler)

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	log.Println("Server started on http://localhost:8080")

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
