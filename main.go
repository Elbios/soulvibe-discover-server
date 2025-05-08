package main

import (
	"log"
	"net/http"
)

func main() {
	cfg := LoadConfig()
	
	LoadTemplates()
	InitializeJobQueue(10, cfg) // Queue buffer size 10

	mux := http.NewServeMux()

	// Serve static files (CSS, JS)
	staticFS := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", staticFS))

	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/api/submit", submitHandler)
	mux.HandleFunc("/api/status/", statusHandler)

	log.Printf("SoulVibe Discover server starting on port %s", cfg.Port)
	log.Printf("Configuration loaded: %+v", cfg) // Be careful with logging sensitive parts of config like passwords in production
	
	// Mask sensitive data in log output if any
	maskedConfig := *cfg
	maskedConfig.SlskPassword = "[REDACTED]"
	maskedConfig.SpotifyEnvVars = []string{"[REDACTED]"} // Mask all Spotify vars
	log.Printf("Sanitized Configuration loaded: Port=%s,  SlskUsername=%s, CliCommandName=%s, TempOutputDir=%s",
        maskedConfig.Port, maskedConfig.SlskUsername, maskedConfig.CliCommandName, maskedConfig.TempOutputDir)


	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}