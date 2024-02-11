# syntax=docker/dockerfile:1

FROM golang:1.21.4

WORKDIR /app

RUN mkdir artifacts

COPY go.mod go.sum ./
RUN go mod download

COPY ./app ./app

# Install yt-dlp
RUN curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp
RUN chmod a+rx /usr/local/bin/yt-dlp  # Make executable
RUN yt-dlp --version

# Install ffmpeg
RUN apt update
RUN apt install -y ffmpeg
RUN ffmpeg -version

RUN CGO_ENABLED=0 GOOS=linux go build -o ./osynetelegramvideodownloader ./app

CMD ["./osynetelegramvideodownloader"]
