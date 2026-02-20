# Stage 1: Build frontend
FROM node:20-alpine AS frontend-builder
WORKDIR /app/frontend
RUN corepack enable && corepack prepare pnpm@latest --activate
COPY frontend/package.json frontend/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile
COPY frontend/ .
RUN pnpm build

# Stage 2: Build backend
FROM golang:1.23-alpine AS backend-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o tracking-server ./cmd/server

# Stage 3: Final image
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=backend-builder /app/tracking-server .
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist
COPY configs/config.toml ./configs/config.toml
EXPOSE 8091
CMD ["./tracking-server", "-c", "config", "-cPath", "./,./configs/"]
