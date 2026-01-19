# ---------------------------------------------
# Base Stage - Dependências
# ---------------------------------------------
FROM golang:1.25.5 AS base
WORKDIR /app

# Cache dos módulos
COPY go.mod go.sum ./
RUN go mod download

# Copia todo o código fonte
COPY . .

# ---------------------------------------------
# Stage de Debug (DEV)
# ---------------------------------------------
FROM base AS dev
# Instala o Delve
RUN go install github.com/go-delve/delve/cmd/dlv@latest

EXPOSE 8080 2345

# COMANDO DE DEBUG OTIMIZADO:
# 1. go build ... ./cmd: Aponta para a pasta onde está o main.go
# 2. -o /tmp/debug-app: Salva o binário fora do volume montado (corrige erro do Windows)
CMD ["sh", "-c", "go build -gcflags='all=-N -l' -o /tmp/debug-app ./cmd && dlv exec --headless --listen=:2345 --api-version=2 --accept-multiclient --continue /tmp/debug-app"]

# ---------------------------------------------
# Stage de Produção (PROD)
# ---------------------------------------------
FROM base AS builder
# Build para produção (também apontando para ./cmd)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server ./cmd

FROM alpine:latest AS prod
WORKDIR /root/
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]