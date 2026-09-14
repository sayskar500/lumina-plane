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
echo "🌐 Starting AI Gateway Server..."
podman run -d --name lumina-server \
  -p 8000:8000 \
  --link lumina-mongo:mongodb \
  -e MONGO_URI="mongodb://mongodb:27017" \
  -e GROQ_API_KEY="your_key_here" \
  lumina-plane-server

echo "✅ Local deployment complete! Try running: lumina health"
