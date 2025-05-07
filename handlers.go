package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var templates *template.Template

func LoadTemplates() {
	var err error
	// Adjust path if your executable is not in the project root
	templates, err = template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatalf("Error loading templates: %v", err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	err := templates.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

type SubmitRequest struct {
	Query string `json:"query"`
}

type SubmitResponse struct {
	JobID string `json:"job_id"`
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SubmitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if req.Query == "" {
		http.Error(w, "Query cannot be empty", http.StatusBadRequest)
		return
	}

	log.Printf("Received submission for query: %s", req.Query)
	job := SubmitJob(req.Query)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SubmitResponse{JobID: job.ID})
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	jobID := filepath.Base(r.URL.Path) // Extracts job_id from /api/status/job_id
	if jobID == "" || jobID == "status" {
		http.Error(w, "Job ID is required", http.StatusBadRequest)
		return
	}

	job, exists := GetJobStatus(jobID)
	if !exists {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(job)
}