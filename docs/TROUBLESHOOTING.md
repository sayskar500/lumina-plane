# 🛠️ Troubleshooting Guide

This guide covers the most common problems encountered when running Lumina-Plane and how to solve them.

## 🚨 Groq AI Connectivity Issues

If you see "AI Provider not configured" or "AI Generation failed," check the following:

### 1. Missing API Key
The server requires a `GROQ_API_KEY` to function.
- **Problem:** You didn't set the environment variable.
- **Solution:** 
  - If using `docker-compose`, add it to the `environment` section of the `server` service.
  - If running locally: `export GROQ_API_KEY="your_key_here" && go run cmd/server/main.go`.

### 2. 401 Unauthorized
- **Problem:** The API key is invalid or expired.
- **Solution:** Visit [console.groq.com](https://console.groq.com) and generate a new API key.

### 3. 429 Too Many Requests
- **Problem:** You have hit the Groq free tier rate limit.
- **Solution:** Wait a few minutes or try a different model via `LUMINA_MODEL` (e.g., `openai/gpt-oss-20b`).

### 4. 400/404 Model Not Found
- **Problem:** The configured model ID is no longer available. Groq rotates its catalog — older IDs like `llama3-8b-8192` are decommissioned.
- **Solution:** List the currently available models and update `ai.default_model` in `config/config.yaml`:
  ```bash
  curl -s https://api.groq.com/openai/v1/models -H "Authorization: Bearer $GROQ_API_KEY" | jq -r '.data[].id'
  ```

---

## 🛢️ MongoDB Connection Issues

### 1. "Could not connect to MongoDB after 5 attempts"
- **Problem:** The MongoDB container hasn't started yet, or the port is blocked.
- **Solution:** 
  - Check if the container is running: `podman ps`.
  - If it's not running, check logs: `podman logs lumina-mongo`.
  - Ensure port `27017` is not being used by a local MongoDB installation.

### 2. Connection Timeout in Cloud
- **Problem:** MongoDB Atlas is rejecting the connection.
- **Solution:** Go to MongoDB Atlas $\rightarrow$ **Network Access** $\rightarrow$ Add IP Address $\rightarrow$ Allow `0.0.0.0/0`.

---

## 🌐 Networking & Ports

### 1. Port Conflicts
Lumina-Plane uses the following ports:
- `8000`: AI Gateway Server
- `3001`: Grafana Dashboard
- `27017`: MongoDB
- `4317/4318`: OTel Collector

**Solution:** If a port is in use, find the process: `lsof -i :<port>` and kill it: `kill -9 <pid>`.

### 2. CLI cannot reach Server
- **Problem:** You are running the CLI in a container but the server is on the host.
- **Solution:** Use `http://host.docker.internal:8000` instead of `localhost:8000` in your CLI settings.
