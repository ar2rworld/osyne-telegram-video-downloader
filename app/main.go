package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"ar2rworld/golang-telegram-video-downloader/app/cleaner"
	"ar2rworld/golang-telegram-video-downloader/app/downloader"
	"ar2rworld/golang-telegram-video-downloader/app/httpclient"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
    bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
    if err != nil {
        panic(err)
    }

    ADMIN_ID, err := strconv.ParseInt(os.Getenv("ADMIN_ID"), 10, 64)
    if err != nil {
        panic(err)
    }

    bot.Debug = true

    updateConfig := tgbotapi.NewUpdate(0)

    updateConfig.Timeout = 30
    updates := bot.GetUpdatesChan(updateConfig)

    //hello message to admin
    helloMessage := tgbotapi.NewMessage(ADMIN_ID, "Hello, boss")
    bot.Send(helloMessage)

    for update := range updates {
        if update.Message == nil {
            continue
        }

        messageText := update.Message.Text

        isInstagramRequest := strings.Contains(messageText, "instagram.com")

        if strings.Contains(messageText, "tiktok.com") ||
        strings.Contains(messageText, "twitter.com") ||
        isInstagramRequest ||
        strings.Contains(messageText, "youtube.com/shorts") {
            log.Println("*** Got request to download video")
            var downloadError error
            messageText = cleaner.CleanUrl(messageText)

            if isInstagramRequest && os.Getenv("INSTAGRAM_COOKIES_STRING") != "" {
                instagramAuthClient, err := httpclient.NewHttpClientFromString(os.Getenv("INSTAGRAM_COOKIES_STRING"))
                if err != nil {
                    log.Println(err)
                    continue
                }
                downloadError = downloader.DownloadVideo(messageText, instagramAuthClient)
            } else if isInstagramRequest {
                downloadError = errors.New("I see that you are trying to share from instagram, but I don't have env var defined")
            } else {
                downloadError = downloader.DownloadVideo(messageText, &http.Client{})
            }
            if downloadError != nil {
                log.Println(downloadError)
                log.Println(update.Message)
                continue
            }
            log.Println("*** Downloaded video without errors")
            videoMessage := tgbotapi.NewVideo(update.Message.Chat.ID, tgbotapi.FilePath("./output"))

            videoMessage.ReplyToMessageID = update.Message.MessageID

        	log.Println("*** Started sending video")
            if _, err := bot.Send(videoMessage); err != nil {
          	    log.Println(err)
        	}
            log.Println("*** Finished sending video")
        } else if messageText == "osyndaisyn ba?" {
            message := tgbotapi.NewMessage(update.Message.Chat.ID, "osyndaymyn")
            bot.Send(message)
        }
    }
}
