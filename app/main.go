package main

import (
	"log"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/ar2rworld/golang-telegram-video-downloader/app/handler"
	"github.com/ar2rworld/golang-telegram-video-downloader/app/httpclient"
	"github.com/ar2rworld/golang-telegram-video-downloader/app/match"
	"github.com/ar2rworld/golang-telegram-video-downloader/app/myerrors"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		panic(err)
	}

	adminID, err := strconv.ParseInt(os.Getenv("ADMIN_ID"), 10, 64)
	if err != nil {
		panic(err)
	}

	bot.Debug = false

	updateConfig := tgbotapi.NewUpdate(0)

	updateConfig.Timeout = 30
	updates := bot.GetUpdatesChan(updateConfig)

	instagramAuthClient, err := httpclient.NewHTTPClientFromString(os.Getenv("INSTAGRAM_COOKIES_STRING"))
	if err != nil {
		log.Println(err)
	}

	// hello message to admin
	helloMessage := tgbotapi.NewMessage(adminID, "Hello, boss")
	sentMessage, err := bot.Send(helloMessage)
	myerrors.CheckTextMessage(&helloMessage, err, &sentMessage)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		messageText := update.Message.Text

		url := match.Match(messageText)

		if url != "" {
			err := handler.VideoMessage(update, url, instagramAuthClient, bot)
			if err != nil {
				log.Println(err)
			}
		} else if messageText == "osyndaisyn ba?" {
			message := tgbotapi.NewMessage(update.Message.Chat.ID, "osyndaymyn")
			sentMessage, err := bot.Send(message)
			myerrors.CheckTextMessage(&message, err, &sentMessage)
		}
	}
}
