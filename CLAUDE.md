# Lumina-Plane: AI-Powered Internal Developer Platform

> Working context for AI assistants and contributors. Feature marketing lives
> in [README.md](README.md); agent architecture lives in [agentic_flow.md](agentic_flow.md).

## What This Is
An AI-native Internal Developer Platform (IDP) built in Go. Developers use a
CLI to query an AI Gateway that augments prompts with versioned, per-project
templates and tracks token usage — the LLMOps core — backed by the platform
engineering stack (IaC, containers, Helm, observability).

## Actual Tech Stack (as implemented)
- **CLI:** Go (`cmd/lumina`) — thin HTTP client to the gateway
- **Backend:** Go `net/http` (`cmd/server`) — the AI Gateway
- **Database:** MongoDB (`prompts`, `token_logs`, `counters` collections)
- **AI:** Groq (OpenAI-compatible API) behind the `internal/gateway.Provider` interface
- **Infra:** Podman/Docker Compose, Terraform (MongoDB Atlas + Render **and Azure: App Service F1 + Cosmos DB free tier**), Helm chart
- **CI/CD:** GitHub Actions → GHCR → App Service (see docs/AZURE.md)
- **Observability:** Prometheus `/metrics`, OpenTelemetry traces → OTel Collector → Grafana

> Note: an early plan called for Python/FastAPI and gRPC. The implementation is
> pure Go over HTTP/JSON; treat FastAPI/gRPC references in older notes as
> historical.

## Current Status
- [x] Phase 1 Foundation — CLI → Server → MongoDB connectivity, compose stack
- [x] Phase 2 AI Gateway — Groq integration, token logging, per-request model routing
- [x] Phase 3 Prompt Registry (partial) — versioned prompts with atomic version counter; prompt evaluation pending
- [x] Phase 4 IaC (partial) — Terraform modules + deploy.sh; GitHub Actions pipeline not yet implemented
- [x] Phase 5 Observability (partial) — Prometheus metrics + provisioned Grafana dashboard + OTel traces; AI RCA tool pending

## Configuration Precedence
1. **Environment variables** (12-factor overrides):
   - Server: `SERVER_PORT`, `MONGO_URI`, `OTEL_ENABLED`, `OTEL_EXPORTER_OTLP_ENDPOINT`, `OTEL_SERVICE_NAME`
   - Secrets: `GROQ_API_KEY`, `INFISICAL_CLIENT_ID`, `INFISICAL_CLIENT_SECRET`
   - CLI: `LUMINA_SERVER_URL`, `LUMINA_MODEL`
2. `config/config.yaml` (defaults) and `config/secrets.yaml` (gitignored)

Never commit secrets: `config/secrets.yaml`, `terraform/terraform.tfvars`, and
Terraform state files are gitignored.

## Model Catalog
Groq rotates model IDs frequently (e.g., `llama3-8b-8192` is decommissioned).
Verify against the live catalog before changing `ai.default_model` in
`config/config.yaml`:

```bash
curl -s https://api.groq.com/openai/v1/models -H "Authorization: Bearer $GROQ_API_KEY" | jq -r '.data[].id'
```

## Local Commands
- `make build` / `make up` / `make down` / `make test`
- `./bin/lumina health | ask "..." | prompt set "<project>" "<template>"`
- Full walkthrough: [docs/LOCAL.md](docs/LOCAL.md)
- **Cloud:** [deploy.sh](deploy.sh) (Render/Atlas) · [scripts/deploy-azure.sh](scripts/deploy-azure.sh) (Azure) · [HF Space](docs/HUGGINGFACE.md) (no card ever)

## Development Guidelines
- Professional Go style: explicit error handling, thin HTTP handlers, context propagation.
- Every cloud resource defined in Terraform; secrets only via environment or Infisical.
- Update `agentic_flow.md` whenever agent behavior changes.
