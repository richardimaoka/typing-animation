package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/richardimaoka/typing-animation/go/server/gitpkg"
)

func writeError(w http.ResponseWriter, statusCode int, err error) {
	body := struct {
		Status  string
		Message string
	}{
		"error",
		err.Error(),
	}
	log.Printf("Returning HTTP error, %s", err)

	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Printf("Writing HTTP error response failed, %s", err)
		return
	}
}

func HandleFiles(w http.ResponseWriter, r *http.Request) {
	// Check path parameters
	orgname := r.PathValue("orgname")
	reponame := r.PathValue("reponame")
	if orgname == "" || reponame == "" {
		writeError(w, http.StatusBadRequest, fmt.Errorf("orgname = '%s', reponame = '%s', but neither allows an empty value", orgname, reponame))
		return
	}

	log.Printf("GET /repos/%s/%s/files called", orgname, reponame)

	// Get git repo
	repo, err := gitpkg.OpenOrClone(orgname, reponame)
	if err != nil {
		log.Printf("Error upon getting git repo, %s", err)
		writeError(w, http.StatusBadRequest, fmt.Errorf("orgname = '%s', reponame = '%s' are supposedly invalid", orgname, reponame))
		return
	}

	files, err := gitpkg.RepoFiles(repo)
	if err != nil {
		log.Printf("Error upon getting git files in the repo, %s", err)
		writeError(w, http.StatusInternalServerError, fmt.Errorf("internal error"))
		return
	}

	// Success
	body := struct {
		Orgname string
		Repo    string
		Files   []string
	}{orgname, reponame, files}

	err = json.NewEncoder(w).Encode(body)
	if err != nil {
		log.Printf("Error upon encoding body, %+v, to json, %s", body, err)
		writeError(w, http.StatusInternalServerError, fmt.Errorf("internal error"))
		return
	}

	

}

func Run() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /repos/{orgname}/{reponame}/files", HandleFiles)
	mux.HandleFunc("GET /", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page!")
	})

	port := 8080
	log.Printf("starting server at http://localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}
