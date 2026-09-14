#!/bin/bash
set -e

echo "🚀 Starting Lumina-Plane Local Deployment..."

# 1. Start MongoDB
echo "📦 Starting MongoDB..."
podman run -d --name lumina-mongo -p 27017:27017 mongo:latest

# 2. Build Server Image
echo "🏗️ Building Server Image..."
podman build -t lumina-plane-server -f Dockerfile.server .

# 3. Run Server
# The image ships with the default config baked in; secrets come from the
# environment (GROQ_API_KEY), which overrides config/secrets.yaml.
echo "🌐 Starting AI Gateway Server..."
podman run -d --name lumina-server \
  -p 8000:8000 \
  --link lumina-mongo:mongodb \
  -e MONGO_URI="mongodb://mongodb:27017" \
  -e GROQ_API_KEY="${GROQ_API_KEY:-}" \
  lumina-plane-server

if [ -z "${GROQ_API_KEY:-}" ]; then
    echo "⚠️  GROQ_API_KEY is not set - the server will start, but AI features will be disabled."
    echo "    Export it and re-run, or pass it when starting the server."
fi

echo "✅ Local deployment complete! Try running: lumina health"
