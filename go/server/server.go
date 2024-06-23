package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

func Run() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /repos/{orgname}/{reponame}/files", HandleFiles)
	mux.HandleFunc("GET /", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page!")
	})

	port := 8080
	log.Printf("starting server at port = %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}

func HandleFiles(w http.ResponseWriter, r *http.Request) {
	// Check path parameters
	var err error
	orgname := r.PathValue("orgname")
	if orgname == "" {
		err = errors.New("orgname is missing")
	}
	reponame := r.PathValue("reponame")
	if reponame == "" {
		if err != nil {
			err = fmt.Errorf("%s, reponame is missing", err)
		} else {
			err = errors.New("orgname is missing")
		}
	}

	// Regardless of success/error, use JSON
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-Type", "application/json")

	// If path parameter errors, then respond with an error
	if err != nil {
		body := struct {
			Status  string
			Message string
		}{"error", err.Error()}

		err := json.NewEncoder(w).Encode(body)
		if err != nil {
			log.Printf("Writing HTTP error response in HandlFiles failed, %s", err)
			return
		}
	}

	// Success
	body := struct {
		Orgname string
		Repo    string
	}{orgname, reponame}
	err = json.NewEncoder(w).Encode(body)
	if err != nil {
		log.Printf("Writing HTTP response in HandlFiles failed, %s", err)
		return
	}
}
