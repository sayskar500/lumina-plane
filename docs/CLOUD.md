# ☁️ Cloud Deployment Guide

This guide explains how to deploy Lumina-Plane to a production-ready cloud environment using the "Free Stack" architecture.

## ⚡ Automated Deployment (deploy.sh)

`./deploy.sh` (from the repo root) automates steps 1–5 below:

1. Discovers your Render workspace ID from the API.
2. Runs `terraform apply` — creating the **Atlas M0 cluster first**, then the
   **Render service with the real `MONGO_URI`** (no manual env var patching).
3. Verifies the service's environment variables and prints the live URL.

**Prerequisites:**
- `terraform`, `jq`, `curl` installed.
- All keys present in `terraform/terraform.tfvars` (gitignored — see the existing file for the expected variables).
- ⚠️ **Render requires a payment method on file to create services via the API — even on the free plan.** Add one at [dashboard.render.com/billing](https://dashboard.render.com/billing) or the apply fails with HTTP 402.
- ⚠️ The Atlas connection string Terraform outputs has **no database credentials**. Create a database user (with `readWrite` on `lumina_plane`) in the Atlas UI and embed `user:password@` in the URI if the app fails to authenticate.

---

## 🏗️ The Architecture
- **API Hosting:** [Render](https://render.com) (Free Tier)
- **Database:** [MongoDB Atlas](https://www.mongodb.com/cloud/atlas) (M0 Free Tier)
- **Secrets Management:** [Infisical](https://infisical.com) (Developer Plan)
- **CI/CD:** GitHub Actions (planned — not yet implemented)
- **Infrastructure as Code:** Terraform

---

## 🚀 Step-by-Step Deployment

### Step 1: Set up MongoDB Atlas (The Data Layer)
1. Create a free account at [MongoDB Atlas](https://www.mongodb.com/cloud/atlas).
2. Create a **Shared Cluster (M0)**.
3. Under **Network Access**, allow access from `0.0.0.0/0`.
4. Under **Database Access**, create a user and copy the **Connection String**.

### Step 2: Set up Infisical (The Security Layer)
1. Create a free account at [Infisical](https://infisical.com).
2. Create a project `lumina-plane` and add secrets: `MONGO_URI`, `GROQ_API_KEY`, `RENDER_API_KEY`.

### Step 3: Set up Render (The Compute Layer)
1. Create a free account at [Render](https://render.com).
2. Create a new **Web Service** connected to your GitHub repo.
3. Select **Docker** runtime and set the **Dockerfile path** to `Dockerfile.server`.
4. Link your environment variables to your Infisical project.

### Step 4: Configure GitHub Actions (The Automation Layer)
Add the following secrets to your GitHub Repository:
- `RENDER_DEPLOY_HOOK`: Render's deploy webhook.
- `RENDER_API_KEY`: Your Render API key.
- `GROQ_API_KEY`: Your Groq key.
- `INFISICAL_KEY`: Your Infisical machine identity key.

### Step 5: Apply Infrastructure as Code (IaC)
The project uses Terraform to manage cloud resources.
1. Initialize: `terraform init`
2. Apply: `terraform apply -auto-approve`

---

## 📈 Verification
Once deployed, test your cloud platform:
1. Update the CLI target to your Render URL.
2. Run: `./bin/lumina health`
3. Run: `./bin/lumina ask "Is my cloud platform running?"`
