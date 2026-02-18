# ==============================================================================
# Build stage
# ==============================================================================
FROM golang:1.25-alpine AS builder

WORKDIR /app

# 依存関係を先にコピーしてキャッシュを活用
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピーしてビルド
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /server ./cmd/main.go

# ==============================================================================
# Runtime stage
# ==============================================================================
FROM gcr.io/distroless/static-debian12

COPY --from=builder /server /server

EXPOSE 8080

ENTRYPOINT ["/server"]
