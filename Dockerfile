# syntax=docker/dockerfile:1

FROM golang:1.22

WORKDIR /app

RUN mkdir artifacts

COPY go.mod go.sum ./
RUN go mod download

COPY ./app ./app

# Install build tools and dependencies for Python
RUN apt update && apt install -y \
    curl \
    wget \
    build-essential \
    libssl-dev \
    zlib1g-dev \
    libbz2-dev \
    libreadline-dev \
    libsqlite3-dev \
    libncursesw5-dev \
    xz-utils \
    tk-dev \
    libxml2-dev \
    libxmlsec1-dev \
    libffi-dev \
    liblzma-dev \
    ffmpeg \
    cron \
    && rm -rf /var/lib/apt/lists/*

# Build Python 3.13.6 from source
RUN curl -O https://www.python.org/ftp/python/3.13.6/Python-3.13.6.tgz && \
    tar -xf Python-3.13.6.tgz && \
    cd Python-3.13.6 && \
    ./configure --enable-optimizations && \
    make -j$(nproc) && \
    make altinstall && \
    cd .. && rm -rf Python-3.13.6*

RUN python3.13 --version && pip3.13 --version

# Install yt-dlp
RUN curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp && \
    chmod a+rx /usr/local/bin/yt-dlp && \
    yt-dlp --version

# Set up yt-dlp auto-update
RUN crontab -l | { cat; echo "0 0 * * * yt-dlp -U"; } | crontab -

# Build Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o ./osynetelegramvideodownloader ./app

CMD ["./osynetelegramvideodownloader"]
