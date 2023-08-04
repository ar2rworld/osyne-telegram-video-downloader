package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"ar2rworld/golang-telegram-video-downloader/downloader"

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
            log.Println(update.Message)
            continue
        }

        if strings.Contains(update.Message.Text, "tiktok.com") ||
        strings.Contains(update.Message.Text, "twitter.com") ||
        strings.Contains(update.Message.Text, "instagram.com") {
            log.Println("*** Got request to download video")
            err := downloader.DownloadVideo(update.Message.Text)
            if err != nil {
                log.Println(err)
                log.Println(update.Message)
                continue
            }
            log.Println("*** Downloaded video without errors")
            videoMessage := tgbotapi.NewVideo(update.Message.Chat.ID, tgbotapi.FilePath("./output"))
        	
            videoMessage.ReplyToMessageID = update.Message.MessageID

        	log.Println("*** Started sending video")
            if _, err := bot.Send(videoMessage); err != nil {
          	    log.Println(err)
                bot.Send(tgbotapi.NewMessage(videoMessage.ChatID, "Had problems sending video"))
        	}
            log.Println("*** Finished sending video")
        }
    }    
}
