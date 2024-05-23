package main

import (
	"log"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/ar2rworld/golang-telegram-video-downloader/app/botservice"
	"github.com/ar2rworld/golang-telegram-video-downloader/app/handler"
	"github.com/ar2rworld/golang-telegram-video-downloader/app/match"
	"github.com/ar2rworld/golang-telegram-video-downloader/app/myerrors"
)

func main() { //nolint: funlen,gocyclo,cyclop
	botAPI, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatalln(err)
	}

	adminID, err := strconv.ParseInt(os.Getenv("ADMIN_ID"), 10, 64)
	if err != nil {
		log.Fatalln(err)
	}

	logChannelID, err := strconv.ParseInt(os.Getenv("LOG_CHANNEL_ID"), 10, 64)
	if err != nil {
		log.Fatalln(err)
	}

	cookiesPath := os.Getenv("COOKIES_PATH")

	botAPI.Debug = false

	updateConfig := tgbotapi.NewUpdate(0)

	updateConfig.Timeout = 30
	updates := botAPI.GetUpdatesChan(updateConfig)

	// hello message to admin
	helloMessage := tgbotapi.NewMessage(adminID, "Hello, boss")
	sentMessage, err := botAPI.Send(helloMessage)
	myerrors.CheckTextMessage(&helloMessage, err, &sentMessage)

	h := handler.NewHandler(botAPI)
	botService := botservice.NewBotService(botAPI, logChannelID)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		messageText := update.Message.Text

		url := match.Match(messageText)

		if url != "" { //nolint: nestif
			err := h.VideoMessage(&update, url, cookiesPath)
			if err != nil {
				log.Println(err)

				err = botService.Log(&update, err)
				if err != nil {
					log.Println(err)
				}

				err = h.ThumbDown(&update)
				if err != nil {
					log.Println("Error while reacting:", err)
				}
			}
		} else if messageText == "osyndaisyn ba?" {
			message := tgbotapi.NewMessage(update.Message.Chat.ID, "osyndaymyn")
			sentMessage, err := botAPI.Send(message)
			myerrors.CheckTextMessage(&message, err, &sentMessage)
		}
	}
}
