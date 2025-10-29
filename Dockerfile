# syntax=docker/dockerfile:1

FROM golang:1.22

WORKDIR /app

RUN mkdir artifacts

COPY go.mod go.sum ./
RUN go mod download

COPY ./app ./app

# Install yt-dlp
RUN curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp && \
    chmod a+rx /usr/local/bin/yt-dlp && \
    yt-dlp --version

# Install ffmpeg
RUN apt update
RUN apt install -y ffmpeg
RUN ffmpeg -version

# Install crontab and add yt-dlp updating cronjob
RUN apt-get -y install cron
RUN crontab -l | { cat; echo "0 0 * * * yt-dlp -U"; } | crontab -

RUN CGO_ENABLED=0 GOOS=linux go build -o ./osynetelegramvideodownloader -ldflags "-X github.com/ar2rworld/golang-telegram-video-downloader/app/downloader.YtdlpPath=yt-dlp" ./app

CMD ["./osynetelegramvideodownloader"]
