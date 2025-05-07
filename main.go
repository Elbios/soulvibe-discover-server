package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	cfg := LoadConfig()
	
	// Validate critical paths from config
	if cfg.CliProjectPath == "" {
		log.Fatal("CLI_PROJECT_PATH environment variable must be set.")
	}
	if cfg.CliWorkingDir == "" {
		log.Fatal("CLI_WORKING_DIR environment variable must be set.")
	}
	if _, err := os.Stat(cfg.CliWorkingDir); os.IsNotExist(err) {
		log.Fatalf("CLI_WORKING_DIR (%s) does not exist.", cfg.CliWorkingDir)
	}
    projectFilePath := filepath.Join(cfg.CliWorkingDir, filepath.Base(cfg.CliProjectPath))
    if filepath.IsAbs(cfg.CliProjectPath) { // if path provided is absolute
        projectFilePath = cfg.CliProjectPath
    }
	if _, err := os.Stat(projectFilePath); os.IsNotExist(err) {
		log.Fatalf("CLI_PROJECT_PATH points to a non-existent file: %s (resolved from %s and %s)", projectFilePath, cfg.CliWorkingDir, cfg.CliProjectPath)
	}


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
	log.Printf("Sanitized Configuration loaded: Port=%s, DotnetExePath=%s, CliProjectPath=%s, CliWorkingDir=%s, SlskUsername=%s, CliCommandName=%s, TempOutputDir=%s",
        maskedConfig.Port, maskedConfig.DotnetExePath, maskedConfig.CliProjectPath, maskedConfig.CliWorkingDir, maskedConfig.SlskUsername, maskedConfig.CliCommandName, maskedConfig.TempOutputDir)


	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}