# Base with dependencies cached
FROM golang:1.25-alpine AS base
WORKDIR /app
ENV GOTOOLCHAIN=auto
RUN apk add --no-cache git ca-certificates build-base
COPY go.mod go.sum ./
# Fix go version to 1.22 if needed (some go commands auto-upgrade it)
RUN sed -i 's/^go 1\.24\..*/go 1.22/' go.mod || true
RUN go mod download
COPY . .
RUN sed -i 's/^go 1\.24\..*/go 1.22/' go.mod || true
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
RUN mkdir -p /app/data/providers && chown -R app:app /app/data
COPY --from=builder /out/app /app/app
USER app
ENV API_PORT=8080
EXPOSE 8080
ENTRYPOINT ["/app/app"]


