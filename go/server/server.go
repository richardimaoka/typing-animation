package server

import (
	"fmt"
	"log"
	"net/http"
)

// func HandleFiles(w http.ResponseWriter, req *http.Request) {
// 	req.URL.
// 		fmt.Fprintf(w, "hello\n")
// }

func Run() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /repositories/{orgname}/{reponame}/files", func(w http.ResponseWriter, r *http.Request) {
		orgname := r.PathValue("orgname")
		reponame := r.PathValue("reponame")
		fmt.Fprintf(w, "orgname = %s, reponame = %s", orgname, reponame)
	})
	mux.HandleFunc("GET /a", func(w http.ResponseWriter, r *http.Request) {
		orgname := r.PathValue("orgname")
		reponame := r.PathValue("reponame")
		fmt.Fprintf(w, "orgname = %s, reponame = %s", orgname, reponame)
	})
	mux.HandleFunc("GET /", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page!")
	})

	port := 8080
	log.Printf("starting server at port = %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}
