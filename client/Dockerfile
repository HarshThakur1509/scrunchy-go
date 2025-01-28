
# Multi-stage build for minimal Vite React.js image
FROM node:23-alpine3.20 AS builder

WORKDIR /app

# Install dependencies efficiently
COPY package*.json ./
RUN npm ci --silent

# Copy project files and build
COPY . .
RUN npm run build

# Final lightweight stage
FROM node:23-alpine3.20

WORKDIR /app

# Install serve for hosting static files
RUN npm install -g serve

# Copy build artifacts from builder stage
COPY --from=builder /app/dist ./dist

# Expose port and start server
EXPOSE 5173
CMD ["serve", "-s", "dist", "-l", "5173"]