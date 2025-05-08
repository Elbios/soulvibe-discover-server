package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

// AppConfig struct definition (assuming it exists elsewhere, example for context)
// type AppConfig struct {
// 	CliProjectPath string   // No longer used for command execution, but might be in struct
// 	CliCommandName string
// 	SlskUsername   string
// 	SlskPassword   string
// 	DotnetExePath  string   // No longer used for command execution, but might be in struct
// 	CliWorkingDir  string
// 	SpotifyEnvVars []string
// }

// TrackInfo struct definition (assuming it exists elsewhere, example for context)
// type TrackInfo struct {
// 	// ... fields for track information
// }

func RunCliCommand(config *AppConfig, query string, outputFilePath string, jobID string) ([]TrackInfo, error) {
	log.Printf("[%s] Starting CLI command for query: %s", jobID, query)

	// Path to the pre-compiled binary
	spotseekBinaryPath := "/app/spotseek"

	// Construct the command arguments for the 'spotseek' binary
	// /app/spotseek <CliCommandName> "<QUERY>" --output-json /tmp/job_id.json -u user -p pass
	args := []string{
		config.CliCommandName, // This was the command after "--" in the original dotnet run
		query,
		"--output-json",
		outputFilePath,
		"-u",
		config.SlskUsername,
		"-p",
		config.SlskPassword,
	}

	cmd := exec.Command(spotseekBinaryPath, args...)
	cmd.Dir = config.CliWorkingDir // Set working directory

	// Prepare environment variables
	cmd.Env = os.Environ() // Inherit parent environment
	cmd.Env = append(cmd.Env, config.SpotifyEnvVars...)

	log.Printf("[%s] Executing: %s %s (in %s)", jobID, spotseekBinaryPath, strings.Join(args, " "), config.CliWorkingDir)
	log.Printf("[%s] With extra ENV: %v", jobID, config.SpotifyEnvVars)

	// Capture stdout/stderr for logging
	var outBuilder, errBuilder strings.Builder
	cmd.Stdout = &outBuilder
	cmd.Stderr = &errBuilder

	startTime := time.Now()
	err := cmd.Run()
	duration := time.Since(startTime)

	log.Printf("[%s] CLI stdout: %s", jobID, outBuilder.String())
	log.Printf("[%s] CLI stderr: %s", jobID, errBuilder.String())
	log.Printf("[%s] CLI execution time: %s", jobID, duration)

	if err != nil {
		errMsg := fmt.Sprintf("CLI command failed: %v. Stderr: %s", err, errBuilder.String())
		log.Printf("[%s] Error: %s", jobID, errMsg)
		// Attempt to clean up output file even on error, if it exists
		_ = os.Remove(outputFilePath)
		return nil, fmt.Errorf(errMsg)
	}

	// Check if output file was created
	if _, statErr := os.Stat(outputFilePath); os.IsNotExist(statErr) {
		errMsg := fmt.Sprintf("CLI command succeeded but output file %s not found. Stderr: %s", outputFilePath, errBuilder.String())
		log.Printf("[%s] Error: %s", jobID, errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	// Read and parse the output file
	fileContent, readErr := os.ReadFile(outputFilePath)
	if readErr != nil {
		errMsg := fmt.Sprintf("failed to read CLI output file %s: %v", outputFilePath, readErr)
		log.Printf("[%s] Error: %s", jobID, errMsg)
		_ = os.Remove(outputFilePath) // Clean up
		return nil, fmt.Errorf(errMsg)
	}

	// Clean up the output file
	removeErr := os.Remove(outputFilePath)
	if removeErr != nil {
		log.Printf("[%s] Warning: failed to remove CLI output file %s: %v", jobID, outputFilePath, removeErr)
	}

	var tracks []TrackInfo
	if jsonErr := json.Unmarshal(fileContent, &tracks); jsonErr != nil {
		errMsg := fmt.Sprintf("failed to parse JSON from CLI output: %v. Content: %s", jsonErr, string(fileContent))
		log.Printf("[%s] Error: %s", jobID, errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	log.Printf("[%s] Successfully processed query, found %d tracks.", jobID, len(tracks))
	return tracks, nil
}