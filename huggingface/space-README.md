---
title: Lumina-Plane AI Gateway
emoji: 🚀
colorFrom: blue
colorTo: purple
sdk: docker
app_port: 8000
pinned: false
---

# Lumina-Plane AI Gateway

This [Hugging Face Space](https://huggingface.co/spaces) runs
[Lumina-Plane](https://github.com/sayskar500/lumina-plane) — an AI-native
Internal Developer Platform written in Go — as a Docker container.

- **AI Gateway:** `/ask` augments prompts with versioned, per-project
  templates from MongoDB and routes to Groq LLMs.
- **Health:** `GET /health` (server + database status)
- **Metrics:** `GET /metrics` (Prometheus)

This Space is synced automatically from the GitHub repository by
[.github/workflows/deploy-huggingface.yml](https://github.com/sayskar500/lumina-plane/blob/main/.github/workflows/deploy-huggingface.yml).

*Secrets (MONGO_URI, GROQ_API_KEY) are injected at runtime as Space secrets —
never stored in this repository.*
