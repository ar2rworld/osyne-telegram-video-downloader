package handler

import (
	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/ar2rworld/golang-telegram-video-downloader/app/downloader"
)

func VideoMessage(update tgbotapi.Update, url string, bot *tgbotapi.BotAPI) error {
	log.Println("*** Got request to download video")

	err := downloader.DownloadVideo(url, &http.Client{})
	if err != nil {
		return err
	}
	log.Println("*** Downloaded video without errors")
	videoMessage := tgbotapi.NewVideo(update.Message.Chat.ID, tgbotapi.FilePath("./output"))

	videoMessage.ReplyToMessageID = update.Message.MessageID

	log.Println("*** Started sending video")
	if _, err := bot.Send(videoMessage); err != nil {
		return err
	}
	log.Println("*** Finished sending video")
	return nil
}
