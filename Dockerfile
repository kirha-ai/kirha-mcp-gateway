FROM node:23-alpine

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy application code
COPY . .

# Expose the port (Smithery will handle the actual port mapping)
EXPOSE 3000

# Start the server in HTTP mode
CMD ["node", "server/index.js"]