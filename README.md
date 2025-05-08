# soulvibe-discover-server

`soulvibe-discover-server` is a Go server with a simple frontend. It runs the `spotseek` binary (from the [slsk-radar](https://github.com/Elbios/slsk-radar) project) as a subprocess to discover music. It queues multiple requests from the frontend to process them sequentially.

## Prerequisites

Before you begin, ensure you have the following:

1.  **Soulseek Client:** Download and install the [official Soulseek client](http://www.slsknet.org/news/node/1).
2.  **Soulseek Account:**
    * Open the Soulseek client.
    * Try to log in with a username and password that does not exist.
    * If the client connects, you now own that account. Note down the username and password.
3.  **Spotify API Credentials:**
    * Go to the [Spotify Developer Dashboard](https://developer.spotify.com/dashboard/) and create an application.
    * Note down the **Client ID** and **Client Secret**.
    * Set the Redirect URI to `http://127.0.0.1:5543/callback`.
4.  **Spotify Refresh Token:**
    * You'll need to perform an OAuth flow to obtain a refresh token. You can use a simple script or a tool like Postman for this. The general flow involves:
        * Redirecting the user to Spotify's authorization URL with your Client ID and requested scopes.
        * After user authorization, Spotify will redirect back to your specified `redirect_uri` with an authorization code.
        * Exchange this authorization code (along with your Client ID and Client Secret) for an access token and a refresh token.
    * **Note down your Spotify Refresh Token.**
5.  **Google AI Studio API Key:**
    * Go to [Google AI Studio](https://aistudio.google.com/).
    * Create an API key from the free tier.
    * **Note down your Google API Key.**

**Summary of Credentials to Note Down:**
* Soulseek Username
* Soulseek Password
* Spotify Client ID
* Spotify Client Secret
* Spotify Refresh Token
* Google API Key

## Dependency: slsk-radar (spotseek CLI)

`soulvibe-discover-server` relies on the `spotseek` CLI tool, which is part of the `slsk-radar` project.

1.  **Clone the `slsk-radar` repository (if not already done):**
    ```bash
    git clone [https://github.com/Elbios/slsk-radar.git](https://github.com/Elbios/slsk-radar.git)
    cd slsk-radar
    ```
2.  **Build and Publish `spotseek`:**
    The `slsk-radar` project contains a Dockerfile (`Source/Dockerfile`) that handles the publishing step. You need to extract the necessary files after building/publishing.
    Alternatively, to build and publish manually (example):
    ```bash
    # Navigate to the slsk-radar project directory
    # Adjust path to the .csproj file as necessary
    dotnet publish Source/Assemblies/Spotify.Slsk.Integration.Cli/Spotify.Slsk.Integration.Cli.csproj -c Release -o ./publish_output --self-contained true -r <your-target-runtime> # e.g., linux-x64
    ```
    After publishing, you will find the `spotseek` executable and `appsettings.json` in the output directory (e.g., `./publish_output`).

3.  **Copy Files to `soulvibe-discover-server`:**
    Copy the `spotseek` executable and the `appsettings.json` file from the `slsk-radar` publish output directory into the root directory of your `soulvibe-discover-server` project.

    * `spotseek` (the executable binary)
    * `appsettings.json`

    *Note: The `slsk-radar` project can also be run directly as a CLI tool without this server:*
    ```bash
    # Example:
    /root/.dotnet/dotnet run --project Source/Assemblies/Spotify.Slsk.Integration.Cli/Spotify.Slsk.Integration.Cli.csproj -- soulseek-radar "lcy bad blood" -u YOUR_SOULSEEK_USER -p YOUR_SOULSEEK_PASS
    ```

## Building `soulvibe-discover-server`

1.  **Build the Go binary:**
    ```bash
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o soulvibe_server main.go clirunner.go handlers.go jobqueue.go config.go
    ```
2.  **Build the Docker image:**
    ```bash
    docker build -t soulvibe-discover .
    ```

## Running `soulvibe-discover-server`

Use the following Docker command to run the server:

```bash
docker run -d --name soulvibe_app \
  -p 8080:8080 \
  -e PORT="8080" \
  -e DOTNET_EXE_PATH="UNUSED" \
  -e CLI_PROJECT_PATH="UNUSED" \
  -e CLI_WORKING_DIR="UNUSED" \
  -e CLI_COMMAND_NAME="./spotseek" \ # Ensure spotseek is executable and in the root
  -e SLSK_USERNAME="<YOUR_SOULSEEK_USERNAME>" \
  -e SLSK_PASSWORD="<YOUR_SOULSEEK_PASSWORD>" \
  -e SPOTIFY_CLIENT_ID="<YOUR_SPOTIFY_CLIENT_ID>" \
  -e GOOGLE_API_KEY="<YOUR_GOOGLE_API_KEY>" \
  -e SPOTIFY_CLIENT_SECRET="<YOUR_SPOTIFY_CLIENT_SECRET>" \
  -e SPOTIFY_REFRESH_TOKEN="<YOUR_SPOTIFY_REFRESH_TOKEN>" \
  soulvibe-discover
```

Important:Replace <YOUR_...> placeholders with your actual credentials.
Ensure the spotseek binary is in the root of the soulvibe-discover-server project directory (where the Dockerfile is) and is executable (chmod +x spotseek). 
The CLI_COMMAND_NAME should point to it, typically ./spotseek if it's in the working directory set in your Dockerfile.

# Accessing the Application

Once the Docker container is running, open your web browser and navigate to:http://localhost:8080

# Usage
In the input field on the web page, type the artist and title of a track you like (e.g., Artist Name Song Title). Do not use a hyphen between the artist and title.Click the "Discover" button.Wait a few minutes for the processing to complete. The slsk-radar subprocess can take time to search and analyze.The results will display Spotify links to recommended songs.

#Stopping the Server
`docker stop soulvibe_app`

`docker rm soulvibe_app`

#Viewing Logs
To view the logs of the running application:

`docker logs -f soulvibe_app`
