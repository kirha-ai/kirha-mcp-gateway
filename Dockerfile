# --- Build stage: compile TypeScript ---
FROM node:20-alpine AS builder
WORKDIR /app

COPY package*.json ./
RUN if [ -f package-lock.json ]; then npm ci; else npm install; fi

COPY tsconfig.json ./
COPY src ./src
COPY config.json ./config.json
RUN npm run build

# --- Runtime stage: smaller, production-only deps ---
FROM node:20-alpine AS runner
ENV NODE_ENV=production
WORKDIR /app

COPY package*.json ./
RUN if [ -f package-lock.json ]; then npm ci --omit=dev; else npm install --omit=dev; fi

COPY --from=builder /app/build ./build
COPY config.json ./config.json

COPY manifest.json ./manifest.json

ENV PORT=3400
EXPOSE 3400

CMD ["node", "build/index.js"]
