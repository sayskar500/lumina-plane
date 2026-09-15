# 🤗 Hugging Face Space Deployment (Free — No Credit Card, Ever)

The most durable **no-credit-card** option for a working professional.
Hugging Face Spaces run Docker containers on the free CPU tier:

| Limit | Free tier |
| :--- | :--- |
| vCPU / RAM | **2 vCPU / 16 GB** |
| Cost | Free, public Space, no card, no trial clock |
| Sleeping | After ~48h of inactivity (next request wakes it) |
| Storage | Ephemeral (fine — state lives in MongoDB Atlas) |

The database stays on **MongoDB Atlas M0** (free, no card) and AI stays on
**Groq** (free tier) — so the entire stack costs nothing and no platform ever
asks for payment details.

## 🚀 One-time setup (~10 minutes)

### 1. Create the Space
1. Sign up at [huggingface.co](https://huggingface.co) (free, no card).
2. **New Space** → name it `lumina-plane` → SDK: **Docker** → **Blank** → Public → Create.

### 2. Add your secrets to the Space
Space → **Settings** → **Variables and secrets** → add two **secrets**:

| Name | Value |
| :--- | :--- |
| `MONGO_URI` | Your Atlas connection string **with credentials** — print it with `terraform -chdir=terraform output -raw mongodb_connection_string` (the `lumina_app` user Terraform created) |
| `GROQ_API_KEY` | Your Groq key (same one as in `terraform/terraform.tfvars`) |

### 3. Create an HF token and connect CI
1. HF → Settings → **Access Tokens** → New token → type **Write**.
2. Add it to this GitHub repo as the secret **`HF_TOKEN`**.

That's it — [.github/workflows/deploy-huggingface.yml](../.github/workflows/deploy-huggingface.yml)
now mirrors the app source to the Space on every push to `main`
(app source only — an explicit allowlist; Terraform state and secrets never
ship to the Space). HF builds the image and starts the container; the Space
README sets `app_port: 8000`, which our server already listens on.

## 🧪 Verify
```bash
curl https://sayskar500-lumina-plane.hf.space/health
LUMINA_SERVER_URL="https://sayskar500-lumina-plane.hf.space" ./bin/lumina ask "Hello HF!"
```
(The first request after a 48h sleep takes ~1 minute while the container wakes.)

## 📝 Notes
- The Space is public (free tier requirement) — safe, because secrets are
  injected at runtime and the source is already open on GitHub.
- Space logs: Space page → **Logs** (build + runtime).
- To rebuild manually after changing secrets, use **Factory rebuild** in the
  Space settings, or just re-run this workflow (`workflow_dispatch`).
