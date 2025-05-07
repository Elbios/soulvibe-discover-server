package main

import (
	"log"
	"os"
)

type AppConfig struct {
	Port                string
	DotnetExePath       string
	CliProjectPath      string
	CliWorkingDir       string
	SlskUsername        string
	SlskPassword        string
	SpotifyEnvVars      []string // Expects "KEY=VALUE" format
	CliCommandName      string   // e.g., "soulseek-radar"
	TempOutputDir       string   // e.g., "/tmp"
}

func LoadConfig() *AppConfig {
	spotifyVars := []string{}
	// Assuming your Spotify env vars are SPOTIFY_ENV_1, SPOTIFY_ENV_2, SPOTIFY_ENV_3
	// You should adjust these to the actual names.
	// The value for the Go server should be "ACTUAL_ENV_VAR_NAME_FOR_CLI=VALUE"
	// Or, more simply, the Go server will read specific env vars and reconstruct them for the CLI.
	// Let's assume the Go server gets SPOTIFY_CLIENT_ID, SPOTIFY_CLIENT_SECRET, etc., directly.
	
	// Example: if your CLI needs FOO=bar and BAZ=qux
	// Env for Go server: SPOTIFY_ENV_VARS="FOO,BAZ"
	// Then Go server reads env FOO and env BAZ
	// For simplicity, let's assume specific known Spotify env vars
	// These are the names of the env vars your .NET CLI expects.
	// The Go server itself will also need these values, so it will read them from its own environment.
	
	// The Go server will read these env vars:
	// SPOTIFY_CLIENT_ID, SPOTIFY_CLIENT_SECRET, SPOTIFY_REFRESH_TOKEN (example names)
	// And pass them to the CLI with the same names.

	// For this example, let's assume the .NET CLI expects these three env vars:
	// `CLI_SPOTIFY_CLIENT_ID`, `CLI_SPOTIFY_CLIENT_SECRET`, `CLI_SPOTIFY_REFRESH_TOKEN`
	// The Go server will read these from its *own* environment (e.g., GO_SPOTIFY_CLIENT_ID)
	// and then set them for the CLI. This is more robust.

	// Let's list the env vars the .NET CLI *expects*:
	cliSpotifyEnvNames := []string{"SPOTIFY_CLIENT_ID", "SPOTIFY_CLIENT_SECRET", "SPOTIFY_REFRESH_TOKEN", "GOOGLE_API_KEY"} // Customize these!

	for _, name := range cliSpotifyEnvNames {
		val := getEnv(name, "") // Go server reads env var intended for CLI
		if val != "" {
			spotifyVars = append(spotifyVars, name+"="+val)
		} else {
			log.Printf("Warning: Environment variable %s for CLI not set for the Go server.", name)
		}
	}


	return &AppConfig{
		Port:                getEnv("PORT", "8080"),
		DotnetExePath:       getEnv("DOTNET_EXE_PATH", "/root/.dotnet/dotnet"),
		CliProjectPath:      getEnv("CLI_PROJECT_PATH", ""), // e.g., /app/cli_project/Source/....csproj - MUST BE SET
		CliWorkingDir:       getEnv("CLI_WORKING_DIR", ""),   // e.g., /app/cli_project/ - MUST BE SET
		SlskUsername:        getEnv("SLSK_USERNAME", ""),    // MUST BE SET
		SlskPassword:        getEnv("SLSK_PASSWORD", ""),    // MUST BE SET
		SpotifyEnvVars:      spotifyVars,
		CliCommandName:      getEnv("CLI_COMMAND_NAME", "soulseek-radar"),
		TempOutputDir:       getEnv("TEMP_OUTPUT_DIR", "/tmp"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	if fallback == "" && (key == "CLI_PROJECT_PATH" || key == "CLI_WORKING_DIR" || key == "SLSK_USERNAME" || key == "SLSK_PASSWORD") {
		log.Fatalf("FATAL: Environment variable %s is required but not set.", key)
	}
	return fallback
}