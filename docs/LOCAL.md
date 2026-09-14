# 💻 Local Development Guide

This guide will help you get Lumina-Plane running on your local machine.

## 🛠️ Prerequisites

Ensure you have the following installed:
- **Go** (v1.22+)
- **Podman** (or Docker)
- **Podman-Compose** (or Docker-Compose)
- **Groq API Key** (Get one for free at [console.groq.com](https://console.groq.com))

## 🚀 Quick Start

### 1. Configure Your Secrets
Lumina-Plane uses a professional YAML-based configuration system.
1. Create the secrets file:
   ```bash
   touch config/secrets.yaml
   ```
2. Add your API key to `config/secrets.yaml`:
   ```yaml
   groq_api_key: "gsk_your_actual_key_here"
   infisical_key: "your_infisical_key_here"
   ```

### 2. Start the Infrastructure
The platform uses a multi-container stack (MongoDB, AI Server, OTel Collector, and Grafana).

```bash
# Start the full stack in the background
podman-compose up -d
```

### 3. Build and Run the CLI
The `lumina` CLI is the primary interface for the platform.

```bash
# Build the binary using the Makefile
make build

# Run the health check
./bin/lumina health
```

## 🕹️ Using the Platform

### Check System Health
Verify that the CLI can talk to the Server, and the Server can talk to MongoDB.
```bash
./bin/lumina health
```
**Expected Output:** `✅ Platform Status: {"status": "healthy", "database": "connected"}`

### Set a Project Prompt Template
The Prompt Registry allows you to define how the AI should behave for a specific project.
```bash
# Usage: lumina prompt set "<project_id>" "<template>"
./bin/lumina prompt set "customer-bot" "You are a professional support agent for a cloud company. Be concise and helpful."
```

### Query the AI Gateway
Ask a question. If a prompt template was set for the project, the AI will follow those instructions.
```bash
# Usage: lumina ask "<prompt>"
./bin/lumina ask "How do I reset my password?"
```

## 📊 Monitoring the Platform
The platform is equipped with a full observability pipeline.
- **Grafana Dashboard:** Open `http://localhost:3001` to see real-time AI performance metrics.
- **Prometheus Metrics:** Access raw metrics at `http://localhost:8000/metrics`.

## 🛠️ Troubleshooting
- **Port 8000 already in use:** Check for other processes (`lsof -i :8000`) and kill them.
- **MongoDB Connection Error:** Ensure the `mongodb` container is healthy (`podman ps`).
- **AI Response Error:** Check `config/secrets.yaml` for a valid API key.
