FROM golang:1.25.5 AS builder
WORKDIR /src

# Cache go mod downloads
COPY go.mod go.sum ./
RUN go env -w GOPROXY=https://proxy.golang.org && go mod download

# Copy full project and build optimized static binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags='-s -w' -o /app/portfolio ./cmd

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/portfolio /app/portfolio
ENV PORT=8080
EXPOSE 8080
USER 1000
CMD ["/app/portfolio"]
