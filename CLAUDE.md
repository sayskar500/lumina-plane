# Lumina-Plane: AI-Powered Internal Developer Platform

Lumina-Plane is a professional-grade project designed to showcase Platform Engineering and AI Platform Engineering skills. It allows developers to self-provision AI-enabled environments through a Go CLI and a FastAPI orchestration layer.

## 🎯 Resume Goals
- **Platform Engineering:** Showcase Self-Service Infra, IaC (Terraform), and GitOps.
- **AI Platform Engineering:** Showcase LLMOps, Model Gateways, and Prompt Registry.
- **Systems Engineering:** Showcase Go, gRPC, and Container Orchestration (Podman/K8s).

## 🛠️ Tech Stack
- **CLI:** Go (Golang)
- **Backend:** Python / FastAPI
- **Database:** MongoDB (Atlas for Cloud / Community for Local)
- **AI Intelligence:** Groq / Google Gemini (via API Gateway)
- **Infrastructure:** Terraform, GitHub Actions, Render/Oracle Cloud
- **Containers:** Podman, Kubernetes (Kind/Minikube)
- **Secrets:** Infisical / HashiCorp Vault

## 🗺️ Implementation Roadmap

### Phase 1: Foundation (Current)
- [ ] Project structure initialization.
- [ ] Basic connectivity: Go CLI $\rightarrow$ FastAPI $\rightarrow$ MongoDB.
- [ ] Local development environment using Podman.

### Phase 2: The AI Gateway (LLMOps)
- [ ] Integration with Groq/Gemini APIs.
- [ ] Centralized LLM Gateway in FastAPI.
- [ ] Token usage tracking and logging in MongoDB.
- [ ] Model routing logic.

### Phase 3: The Prompt Registry
- [ ] Version-controlled prompt store in MongoDB.
- [ ] API for retrieving "Production" vs "Staging" prompts.
- [ ] Prompt evaluation tool.

### Phase 4: Cloud Automation (IaC)
- [ ] Terraform modules for MongoDB Atlas and Render.
- [ ] GitHub Actions pipeline for automated deployment.
- [ ] Implementation of GitOps flow.

### Phase 5: Observability & SRE
- [ ] Prometheus metrics for API latency and token cost.
- [ ] Grafana dashboard for platform health.
- [ ] Automated Root Cause Analysis (RCA) tool using AI.

## 🚦 Development Guidelines
- Match the style of professional Go and Python code.
- Every single cloud resource must be defined in Terraform.
- Use gRPC for internal communication between Go and Python where performance is critical.
- All secrets must be managed via Infisical/Vault, never in `.env` files.
