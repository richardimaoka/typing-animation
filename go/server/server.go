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

func HandleRepoFiles(w http.ResponseWriter, r *http.Request) {
	// Check path parameters
	orgname := r.PathValue("orgname")
	reponame := r.PathValue("reponame")
	if orgname == "" || reponame == "" {
		writeError(w, http.StatusBadRequest, fmt.Errorf("orgname = '%s', reponame = '%s', but neither allows an empty value", orgname, reponame))
		return
	}

	// Path parameter checks passed
	log.Printf("GET /repos/%s/%s/files called", orgname, reponame)

	// Get repo files
	files, err := gitpkg.RepoFiles(orgname, reponame)
	if err != nil {
		log.Printf("Error upon getting git files in the repo, %s", err)
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	// Success
	body := struct {
		Orgname string   `json:"orgname"`
		Repo    string   `json:"repo"`
		Files   []string `json:"files"`
	}{orgname, reponame, files}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(body)
	if err != nil {
		log.Printf("Error upon encoding body, %+v, to json, %s", body, err)
		writeError(w, http.StatusInternalServerError, fmt.Errorf("internal error"))
		return
	}
}

func HandleSingleFile(w http.ResponseWriter, r *http.Request) {
	// Check path parameters
	orgname := r.PathValue("orgname")
	reponame := r.PathValue("reponame")
	filepath := r.PathValue("filepath")
	if orgname == "" || reponame == "" || filepath == "" {
		writeError(
			w,
			http.StatusBadRequest,
			fmt.Errorf("orgname = '%s', reponame = '%s', filepath = '%s', but neither allows an empty value", orgname, reponame, filepath),
		)
		return
	}

	commitHash := r.URL.Query().Get("commit")

	// Path parameter checks passed
	log.Printf("GET /repos/%s/%s/files/%s called", orgname, reponame, filepath)

	// Get git repo, then get repo files
	contents, err := gitpkg.RepoFileContents(orgname, reponame, filepath, commitHash)
	if err != nil {
		log.Printf("Error upon getting git file in the repo, %s", err)
		writeError(w, http.StatusInternalServerError, fmt.Errorf("internal error"))
		return
	}

	// Success
	body := struct {
		Orgname  string `json:"orgname"`
		Repo     string `json:"repo"`
		Contents string `json:"contents"`
	}{orgname, reponame, contents}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(body)
	if err != nil {
		log.Printf("Error upon encoding body, %+v, to json, %s", body, err)
		writeError(w, http.StatusInternalServerError, fmt.Errorf("internal error"))
		return
	}
}

func Run() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /repos/{orgname}/{reponame}/files", HandleRepoFiles)
	mux.HandleFunc("GET /repos/{orgname}/{reponame}/files/{filepath...}", HandleSingleFile)

	port := 8080
	log.Printf("starting server at http://localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}
