# syntax=docker/dockerfile:1
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/ ./cmd/
COPY internal/ ./internal/
RUN CGO_ENABLED=0 go build -o /bright-data-mcp ./cmd/api

FROM alpine:3.19
RUN apk --no-cache add ca-certificates curl
RUN addgroup -g 1000 -S appuser && \
    adduser -u 1000 -S appuser -G appuser -s /sbin/nologin
COPY --from=builder /bright-data-mcp /bright-data-mcp
WORKDIR /
EXPOSE 8011
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8011/health || exit 1
USER appuser
CMD ["/bright-data-mcp"]
