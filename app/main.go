package main

import (
	"errors"
	"log"
	"os"
	"strconv"
	"strings"

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
	instagramCookiesPath := os.Getenv("INSTAGRAM_COOKIES_PATH")
	googleCookiesPath := os.Getenv("GOOGLE_COOKIES_PATH")

	botAPI.Debug = false

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := botAPI.GetUpdatesChan(updateConfig)

	// hello message to admin
	helloMessage := tgbotapi.NewMessage(adminID, "Hello, boss")
	sentMessage, err := botAPI.Send(helloMessage)
	myerrors.CheckTextMessage(&helloMessage, err, &sentMessage)

	botService := botservice.NewBotService(botAPI, logChannelID)
	h := handler.NewHandler(botAPI, botService, cookiesPath, instagramCookiesPath, googleCookiesPath)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		messageText := update.Message.Text

		if len(update.Message.Entities) > 0 && update.Message.ReplyToMessage != nil && strings.Contains(messageText, botAPI.Self.UserName) {
			err = h.HandleMentionMessage(&update)
			if err != nil && errors.Is(err, handler.ErrNoURLFound) {
				err = h.Whaat(&update)
				h.HandleError(&update, err)
			} else if err != nil {
				h.HandleError(&update, err)
			}
			continue
		}

		// Inside the main loop where you handle updates
		if update.Message.From.ID == adminID && update.Message.Document != nil {
			err := h.HandleAdminMessage(&update)
			h.HandleError(&update, err)
			continue
		}

		url := match.Match(messageText)
		if url != "" {
			err := h.VideoMessage(&update, url)
			if err != nil {
				h.HandleError(&update, err)
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
