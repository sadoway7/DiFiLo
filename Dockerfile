# syntax=docker/dockerfile:1

# ---- build stage: compile the Go binary for Linux ----
FROM golang:1.25-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
COPY cmd/ internal/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /difilo ./cmd/difilo

# ---- runtime stage: binary + mirror, nothing else ----
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /difilo /app/difilo
COPY mirror /app/mirror
# Data directory for the SQLite database (mount as volume to persist across deploys)
RUN mkdir -p /app/data
VOLUME /app/data
# Symlink the database into the data volume so it survives container recreation
ENV DIFILO_DB=/app/data/difilo.db
# Bind 0.0.0.0 so the published port reaches the server inside the container.
EXPOSE 8000
ENTRYPOINT ["/app/difilo", "--mirror", "/app/mirror", "--host", "0.0.0.0", "--port", "8000", "--db", "/app/data/difilo.db"]
