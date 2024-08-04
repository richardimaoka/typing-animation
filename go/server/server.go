package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/richardimaoka/typing-animation/go/gitpkg"
)

func writeErrorJson(w http.ResponseWriter, statusCode int, err error) {
	body := struct {
		Status  string `json:"status"`
		Message string `json:"message"`
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

func HandleGET_Repo(w http.ResponseWriter, r *http.Request) {
	// Check path parameters
	orgname := r.PathValue("orgname")
	reponame := r.PathValue("reponame")
	if orgname == "" || reponame == "" {
		writeErrorJson(w, http.StatusBadRequest, fmt.Errorf("orgname = '%s', reponame = '%s', but neither allows an empty value", orgname, reponame))
		return
	}

	// Path parameter checks passed
	log.Printf("GET /repos/%s/%s called", orgname, reponame)

	// Start git clone in goroutine
	_, err := gitpkg.Open(orgname, reponame)
	if err != nil {
		log.Printf("Repo is not ready, %s", err)
		writeErrorJson(w, http.StatusInternalServerError, err)
		return
	}

	// Success
	body := struct {
		Orgname string `json:"orgname"`
		Repo    string `json:"repo"`
		Status  string `json:"status"`
	}{orgname, reponame, "ready"}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(body)
	if err != nil {
		log.Printf("Error upon encoding body, %+v, to json, %s", body, err)
		writeErrorJson(w, http.StatusInternalServerError, fmt.Errorf("internal error"))
		return
	}
}

func HandlePOST_Repo(w http.ResponseWriter, r *http.Request) {
	// Check path parameters
	orgname := r.PathValue("orgname")
	reponame := r.PathValue("reponame")
	if orgname == "" || reponame == "" {
		writeErrorJson(w, http.StatusBadRequest, fmt.Errorf("orgname = '%s', reponame = '%s', but neither allows an empty value", orgname, reponame))
		return
	}

	// Path parameter checks passed
	log.Printf("POST /repos/%s/%s called", orgname, reponame)

	// Start git clone in goroutine
	var err error

	// Success
	body := struct {
		Orgname string `json:"orgname"`
		Repo    string `json:"repo"`
	}{orgname, reponame}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(body)
	if err != nil {
		log.Printf("Error upon encoding body, %+v, to json, %s", body, err)
		writeErrorJson(w, http.StatusInternalServerError, fmt.Errorf("internal error"))
		return
	}
}

func HandleRepoFiles(w http.ResponseWriter, r *http.Request) {
	// Check path parameters
	orgname := r.PathValue("orgname")
	reponame := r.PathValue("reponame")
	if orgname == "" || reponame == "" {
		writeErrorJson(w, http.StatusBadRequest, fmt.Errorf("orgname = '%s', reponame = '%s', but neither allows an empty value", orgname, reponame))
		return
	}

	// Path parameter checks passed
	log.Printf("GET /repos/%s/%s/files called", orgname, reponame)

	// Get repo files
	files, err := gitpkg.RepoFiles(orgname, reponame)
	if err != nil {
		log.Printf("Error upon getting git files in the repo, %s", err)
		writeErrorJson(w, http.StatusInternalServerError, err)
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
		writeErrorJson(w, http.StatusInternalServerError, fmt.Errorf("internal error"))
		return
	}
}

func HandleSingleFile(w http.ResponseWriter, r *http.Request) {
	// Check path parameters
	orgname := r.PathValue("orgname")
	reponame := r.PathValue("reponame")
	filepath := r.PathValue("filepath")
	if orgname == "" || reponame == "" || filepath == "" {
		writeErrorJson(
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
	commits, err := gitpkg.CommitsForFile(orgname, reponame, filepath)
	if err != nil {
		log.Printf("Error upon getting git file in the repo, %s", err)
		writeErrorJson(w, http.StatusInternalServerError, fmt.Errorf("internal error"))
		return
	}

	type CommitData struct {
		Hash         string `json:"hash"`
		ShortHash    string `json:"shortHash"`
		Message      string `json:"message"`
		ShortMessage string `json:"shortMessage"`
	}
	var commitDataSlice []CommitData
	for _, c := range commits {
		messageInRunes := []rune(c.Message)
		var shortMessage string
		if len(messageInRunes) > 20 {
			shortMessage = string(messageInRunes[:20]) + "..."
		} else {
			shortMessage = string(messageInRunes)
		}

		hash := c.Hash.String()
		data := CommitData{
			Hash:         hash,
			ShortHash:    string([]rune(hash)[:7]),
			Message:      c.Message,
			ShortMessage: shortMessage,
		}
		commitDataSlice = append(commitDataSlice, data)
	}

	// Success
	body := struct {
		Orgname string       `json:"orgname"`
		Repo    string       `json:"repo"`
		Commits []CommitData `json:"commits"`
	}{orgname, reponame, commitDataSlice}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(body)
	if err != nil {
		log.Printf("Error upon encoding body, %+v, to json, %s", body, err)
		writeErrorJson(w, http.StatusInternalServerError, fmt.Errorf("internal error"))
		return
	}
}

func HandleSingleCommit(w http.ResponseWriter, r *http.Request) {
	// Check path parameters
	orgname := r.PathValue("orgname")
	reponame := r.PathValue("reponame")
	filepath := r.PathValue("filepath")
	if orgname == "" || reponame == "" || filepath == "" {
		writeErrorJson(
			w,
			http.StatusBadRequest,
			fmt.Errorf("orgname = '%s', reponame = '%s', filepath = '%s', but neither allows an empty value", orgname, reponame, filepath),
		)
		return
	}

	commitHash := r.URL.Query().Get("commit")
	if commitHash == "" {
		writeErrorJson(
			w,
			http.StatusBadRequest,
			fmt.Errorf("query parameter commit is missing"),
		)
		return
	}

	// Path parameter checks passed
	log.Printf("GET /repos/%s/%s/files/%s called", orgname, reponame, filepath)

	// Get git repo, then get repo files
	contents, err := gitpkg.RepoFileContents(orgname, reponame, filepath, commitHash)
	if err != nil {
		log.Printf("Error upon getting git file in the repo, %s", err)
		writeErrorJson(w, http.StatusInternalServerError, fmt.Errorf("internal error"))
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
		writeErrorJson(w, http.StatusInternalServerError, fmt.Errorf("internal error"))
		return
	}
}

func Run() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /{orgname}/{reponame}", HandlePOST_Repo)
	mux.HandleFunc("GET /{orgname}/{reponame}", HandleGET_Repo)
	mux.HandleFunc("GET /{orgname}/{reponame}/files", HandleRepoFiles)
	mux.HandleFunc("GET /{orgname}/{reponame}/files/{filepath...}", HandleSingleFile)

	// mux.HandleFunc("GET /repos/{orgname}/{reponame}/", HandleRepoFiles)
	// mux.HandleFunc("GET /repos/{orgname}/{reponame}/branches", HandleRepoFiles)
	// mux.HandleFunc("GET /repos/{orgname}/{reponame}/branches", HandleRepoFiles)
	// mux.HandleFunc("GET /repos/{orgname}/{reponame}/branches", HandleRepoFiles)
	// mux.HandleFunc("GET /repos/{orgname}/{reponame}/branches", HandleRepoFiles)
	// mux.HandleFunc("GET /repos/{orgname}/{reponame}/branches", HandleRepoFiles)

	port := 8080
	log.Printf("starting server at http://localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}
