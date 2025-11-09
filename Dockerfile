# syntax=docker/dockerfile:1

# Stage 1: build
FROM golang:1.22-alpine AS builder

WORKDIR /build

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY ./app ./app

RUN CGO_ENABLED=0 GOOS=linux go build -o ./osynetelegramvideodownloader ./app

# Stage 2: Runtime
FROM python:3.11-slim

WORKDIR /app

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        ffmpeg \
        curl \
        ca-certificates \
        cron && \
    rm -rf /var/lib/apt/lists/*

# Install yt-dlp
RUN curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp && \
    chmod a+rx /usr/local/bin/yt-dlp && \
    yt-dlp --version

# Install crontab and add yt-dlp updating cronjob
RUN crontab -l | { cat; echo "0 0 * * * yt-dlp -U"; } | crontab -

COPY --from=builder /build/osynetelegramvideodownloader /app/

COPY docker-entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/docker-entrypoint.sh

ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["./osynetelegramvideodownloader"]
