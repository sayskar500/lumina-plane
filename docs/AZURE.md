# ☁️ Azure Deployment Guide (Free — No Credit Card)

Lumina-Plane deploys to **Azure free tiers** using Infrastructure as Code:

| Component | Azure Service | Free tier |
| :--- | :--- | :--- |
| AI Gateway (container) | **App Service** (Linux, F1) | Free forever — 1 GB RAM, sleeps after 20 min idle |
| MongoDB | **Cosmos DB** (MongoDB API, v7.0) | Free tier — 1000 RU/s + 25 GB **forever** |
| Container registry | **GHCR** (GitHub) | Free for public repos |
| CI/CD | **GitHub Actions** | Free for public repos |

**No credit card is charged and nothing can incur cost:** the only throughput
allocation is a shared-throughput database at exactly the free 1000 RU/s, and
App Service F1 is a free SKU.

## ⚠️ About the "no credit card" requirement

- **AWS requires a card** at signup (identity verification), even for its free
  plan — so it does not satisfy a no-card constraint.
- **Azure** requires a card for the *standard free account*, **but not** for:
  - **Azure for Students** — $100 credit + free services, verified with a
    school email, **no card**: https://azure.microsoft.com/free/students
  - **GitHub Student Developer Pack** → includes the Azure student offer:
    https://education.github.com/pack
- If you are not a student, the standard Azure free account needs a
  credit/debit card purely for identity verification ($0 charge unless you
  explicitly upgrade).

## 🚀 Step-by-step

### 1. Get an Azure subscription (no card)
Sign up for **Azure for Students** (or via the GitHub Student Pack). You end
up with a subscription named `Azure for Students` — no payment step.

### 2. Install the Azure CLI and log in
```bash
brew install azure-cli
az login          # opens a browser; pick your Azure for Students subscription
az account show   # verify the subscription name
```

### 3. Publish the container image (CI/CD)
Push to `main` — [.github/workflows/docker-publish.yml](../.github/workflows/docker-publish.yml)
builds the image and publishes it to `ghcr.io/sayskar500/lumina-plane:latest`.

**One-time:** after the first successful run, make the package public so
App Service can pull it anonymously:
GitHub → your repo → **Packages** → `lumina-plane` → **Package settings** →
**Change visibility** → **Public**.

### 4. Deploy with Terraform
```bash
./scripts/deploy-azure.sh
```
The script checks your Azure login, reuses the Groq key from
`terraform/terraform.tfvars`, and creates everything:
resource group → Cosmos DB (free tier) → App Service (F1) with
`MONGO_URI` + `GROQ_API_KEY` wired as app settings.

### 5. (Optional) Auto-redeploy on new images
App Service needs a nudge when a new `:latest` image is pushed:
```bash
az webapp deployment container config --enable-cd true \
  --name <web_app_name> --resource-group <resource_group> --query CI_CD_URL -o tsv
```
Store that URL as the GitHub repo secret `AZURE_CD_WEBHOOK_URL` — the CI
workflow will trigger a redeploy after each image push.

## 🧪 Verify
```bash
curl https://<web_app_url>/health     # first call may take ~1 min (F1 cold start)
LUMINA_SERVER_URL="https://<web_app_url>" ./bin/lumina ask "Hello Azure!"
```

## 🧰 Operations
```bash
az webapp log tail --name <web_app_name> --resource-group <resource_group>
terraform -chdir=terraform/azure output web_app_url
```

## 💸 Cost safety
- App Service F1: free SKU, cannot cost money.
- Cosmos DB free tier: 1000 RU/s + 25 GB free; our single shared-throughput
  database uses exactly 1000 RU/s and collections inherit it.
- Nothing here consumes the $100 student credit. If you ever add paid
  resources, set a spending alert in the portal (Cost Management → Budgets).

## 🗺️ Multi-cloud note
The Render + MongoDB Atlas deployment (`deploy.sh`) remains live and free —
you can keep both as a multi-cloud setup (great portfolio story) or tear
Render down with `terraform destroy` in `terraform/`.
