package handler

import (
	"log"
	"net/http"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/ar2rworld/golang-telegram-video-downloader/app/downloader"
	"github.com/ar2rworld/golang-telegram-video-downloader/app/myerrors"
)

func VideoMessage(update tgbotapi.Update, url string, instagramAuthClient *http.Client, bot *tgbotapi.BotAPI) error {
	isInstagramRequest := strings.Contains(url, "instagram.com")

	log.Println("*** Got request to download video")
	var downloadError error

	switch {
	case isInstagramRequest && os.Getenv("INSTAGRAM_COOKIES_STRING") != "":
		downloadError = downloader.DownloadVideo(url, instagramAuthClient)
	case isInstagramRequest:
		downloadError = myerrors.NewMissingEnvVariableError("INSTAGRAM_COOKIES_STRING")
	default:
		downloadError = downloader.DownloadVideo(url, &http.Client{})
	}
	if downloadError != nil {
		return downloadError
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
