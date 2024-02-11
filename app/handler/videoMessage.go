package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/wader/goutubedl"

	"github.com/ar2rworld/golang-telegram-video-downloader/app/downloader"
	"github.com/ar2rworld/golang-telegram-video-downloader/app/match"
)

const Duration = 10

func VideoMessage(update tgbotapi.Update, url string, bot *tgbotapi.BotAPI) error {
	remove := []string{}
	defer func() {
		for _, fn := range remove {
			err := os.Remove(fn)
			log.Println("*** Removed file: ", fn, "error:", err)
		}
	}()

	log.Println("*** Got request to download video")

	opts := goutubedl.Options{HTTPClient: &http.Client{}, DebugLog: log.Default()}
	if match.Youtube(url) != "" {
		opts.DownloadSections = fmt.Sprintf("*0:0-0:%d", Duration)
		log.Printf("*** Downloading video from Youtube %s\n", opts.DownloadSections)
	}

	// save video with unique names and put them to /tmp, delete after sending or if error

	fileName, err := downloader.DownloadVideo(url, opts)
	if err != nil {
		return err
	}
	log.Println("*** Downloaded video without errors")
	videoMessage := tgbotapi.NewVideo(update.Message.Chat.ID, tgbotapi.FilePath("./"+fileName))

	videoMessage.ReplyToMessageID = update.Message.MessageID

	log.Println("*** Started sending video")
	m, err := bot.Send(videoMessage)
	if err != nil {
		return err
	}
	log.Println(m.Video.FileName, m.Video.MimeType, "duration:", m.Video.Duration, "size:", m.Video.FileSize)
	log.Println("*** Finished sending video")
	return nil
}
