FROM node:23-alpine

WORKDIR /app

# Copy package files
COPY package*.json ./
COPY tsconfig.json ./

# Install dependencies
RUN npm install

# Copy source code
COPY src/ ./src/

# Build the application
RUN npm run build

# Set PORT environment variable to ensure HTTP mode
ENV PORT=3000

# Expose the port (Smithery will handle the actual port mapping)
EXPOSE 3000

# Start the server in HTTP mode
CMD ["node", "dist/index.js"]