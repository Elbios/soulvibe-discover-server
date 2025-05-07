# Stage 1: Build the Go application (optional, can be done outside)
# If you prefer to build Go app inside Docker:
# FROM golang:1.21-alpine as builder
# WORKDIR /app
# COPY go.mod go.sum ./
# RUN go mod download
# COPY . .
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/soulvibe_server .

# Stage 2: Create the final image with .NET runtime
# Using aspnet runtime as it's smaller than SDK but includes what's needed for `dotnet run` essentially.
# If `dotnet run` absolutely needs SDK tools not in runtime, use mcr.microsoft.com/dotnet/sdk:8.0-alpine
FROM mcr.microsoft.com/dotnet/sdk:8.0-alpine AS final 
# sdk image is larger but ensures `dotnet run` works fully. 
# For a smaller image, if your .NET app is published self-contained or as a framework-dependent dll,
# you could use mcr.microsoft.com/dotnet/aspnet:8.0-alpine or runtime:8.0-alpine

WORKDIR /app

# Copy the pre-compiled Go binary (compile it on your WSL: GOOS=linux GOARCH=amd64 go build -o soulvibe_server main.go clirunner.go handlers.go jobqueue.go config.go)
COPY soulvibe_server /app/soulvibe_server

# Copy frontend assets
COPY templates/ /app/templates/
COPY static/ /app/static/

# Expose the port the Go app will listen on
EXPOSE 8080 

# Environment variables for the Go server will be set during `docker run`
# Example: PORT, DOTNET_EXE_PATH, CLI_PROJECT_PATH, CLI_WORKING_DIR, SLSK_USERNAME, SLSK_PASSWORD, SPOTIFY_...

# Entrypoint
CMD ["/app/soulvibe_server"]