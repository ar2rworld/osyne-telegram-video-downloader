module github.com/ar2rworld/golang-telegram-video-downloader

go 1.21

require (
	github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1
	github.com/wader/goutubedl v0.0.0-20240207160746-8b34407df2f3
)

replace github.com/go-telegram-bot-api/telegram-bot-api/v5 => ./telegram-bot-api/
