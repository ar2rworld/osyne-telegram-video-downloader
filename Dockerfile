# syntax=docker/dockerfile:1

FROM golang:1.22

WORKDIR /app

RUN mkdir artifacts

COPY go.mod go.sum ./
RUN go mod download

COPY ./app ./app

# Install dependencies for Python and general tools
RUN apt update && apt install -y \
    software-properties-common \
    curl \
    gnupg \
    lsb-release

# Add deadsnakes PPA for Python 3.13
RUN add-apt-repository ppa:deadsnakes/ppa && apt update && \
    apt install -y python3.13 python3.13-venv python3.13-dev python3-pip && \
    python3.13 --version && pip3 --version

# Install yt-dlp
RUN curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp
RUN chmod a+rx /usr/local/bin/yt-dlp  # Make executable
RUN yt-dlp --version

# Install ffmpeg
RUN apt install -y ffmpeg
RUN ffmpeg -version

# Install crontab and add yt-dlp updating cronjob
RUN apt-get -y install cron
RUN crontab -l | { cat; echo "0 0 * * * yt-dlp -U"; } | crontab -

# Build Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o ./osynetelegramvideodownloader ./app

CMD ["./osynetelegramvideodownloader"]
