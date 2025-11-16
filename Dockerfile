# Base with dependencies cached
FROM golang:1.22-alpine AS base
WORKDIR /app
RUN apk add --no-cache git ca-certificates build-base
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV GOTOOLCHAIN=auto
RUN go mod tidy

# Builder for production binary
FROM base AS builder
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/app ./cmd/api

# Development stage with hot-reload (air)
FROM base AS dev
RUN go install github.com/cosmtrek/air@v1.49.0
EXPOSE 8080
CMD ["air", "-c", ".air.toml"]

# Runtime image
FROM alpine:3.20 AS runtime
RUN addgroup -S app && adduser -S app -G app
WORKDIR /app
COPY --from=builder /out/app /app/app
USER app
ENV API_PORT=8080
EXPOSE 8080
ENTRYPOINT ["/app/app"]


