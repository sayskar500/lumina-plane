# Lumina-Plane: AI-Powered Internal Developer Platform

Lumina-Plane is an AI-native Internal Developer Platform (IDP) designed to showcase advanced Platform Engineering and LLMOps skills.

## 🌟 High-Impact Features for Resume
- **Self-Service Control Plane:** A Go-based CLI that allows developers to interact with AI infrastructure without manual ticket requests.
- **AI Gateway (LLMOps):** A high-performance Go backend that manages LLM API routing, token tracking, and cost optimization.
- **Versioned Prompt Registry:** Implemented a system to version-control prompts in MongoDB, enabling "canary" prompt deployments and rollbacks.
- **Infrastructure as Code (IaC):** Entire cloud stack (MongoDB Atlas, Render) defined using Terraform.
- **Kubernetes Orchestration:** Full Helm chart implementation for scalable deployment to K8s clusters.
- **Observability & SRE:** integrated Prometheus metrics for real-time tracking of AI latency, token usage, and system health.
- **CI/CD GitOps:** Automated deployment pipeline using GitHub Actions for a full "Code $\rightarrow$ Cloud" flow.

## 🛠️ Technical Architecture
- **Languages:** Go (CLI & Server), Bash.
- **Backend:** Go (`net/http`), MongoDB (Managed Atlas).
- **AI Intelligence:** Groq / Google Gemini.
- **Infra:** Terraform, Render, Podman, Kubernetes (Helm).
- **Observability:** Prometheus, Grafana.

## 📖 Documentation
Detailed guides are available in the `/docs` directory:
- [**Local Development Guide**](docs/LOCAL.md) - How to run the platform on your machine.
- [**Cloud Deployment Guide**](docs/CLOUD.md) - How to deploy to the free cloud stack.

## 🚀 Quick Start (Local)
1. Start the infrastructure:
   ```bash
   ./scripts/deploy-local.sh
   ```
2. Build the CLI:
   ```bash
   cd cmd/lumina && go build -o lumina main.go
   ```
3. Check health:
   ```bash
   ./lumina health
   ```

## 📈 Resume Bullet Points (Use these!)
- "Architected an AI-powered Internal Developer Platform (IDP) using Go and MongoDB, reducing infrastructure provisioning time by implementing a self-service CLI."
- "Engineered a centralized LLM Gateway in Go, implementing token tracking and model routing to optimize AI operational costs."
- "Developed a versioned Prompt Registry in MongoDB, enabling prompt lifecycle management and preventing regressions in AI behavior."
- "Automated cloud resource provisioning for MongoDB Atlas and Render using Terraform, implementing a full GitOps CI/CD pipeline via GitHub Actions."
- "Designed and implemented professional Helm charts for Kubernetes deployment, ensuring scalable and consistent application orchestration."
- "Implemented system-wide observability using Prometheus, tracking critical AI metrics like token-per-second and request latency."
