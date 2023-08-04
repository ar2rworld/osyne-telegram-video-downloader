# syntax=docker/dockerfile:1

FROM golang:1.19

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd ./cmd
COPY ./downloader ./downloader
COPY ./httpclient ./httpclient

RUN CGO_ENABLED=0 GOOS=linux go build -o /osynetelegramvideodownloader

CMD ["/osynetelegramvideodownloader"]