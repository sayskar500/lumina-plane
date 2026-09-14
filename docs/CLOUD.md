# ☁️ Cloud Deployment Guide

This guide explains how to deploy Lumina-Plane to a production-ready cloud environment using the "Free Stack" architecture.

## 🏗️ The Architecture
- **API Hosting:** [Render](https://render.com) (Free Tier)
- **Database:** [MongoDB Atlas](https://www.mongodb.com/cloud/atlas) (M0 Free Tier)
- **Secrets Management:** [Infisical](https://infisical.com) (Developer Plan)
- **CI/CD:** GitHub Actions
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
