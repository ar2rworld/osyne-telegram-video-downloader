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

const Duration = 30

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
	isYoutubeVideo := match.Youtube(url) != ""
	if isYoutubeVideo {
		sections, err := parse(update.Message.Text)
		if err != nil {
			log.Println("*** Parsed video sections")
			sections = fmt.Sprintf("*0:0-0:%d", Duration)
		}
		opts.DownloadSections = sections
		log.Printf("*** Downloading video from Youtube %s\n", opts.DownloadSections)
	}

	fileName, err := downloader.DownloadVideo(url, opts)
	if err != nil {
		return err
	}
	remove = append(remove, fileName)
	log.Println("*** Downloaded video without errors")

	if isYoutubeVideo {
		fileName, err = downloader.Convert(fileName)
		if err != nil {
			return err
		}
		remove = append(remove, fileName)
		log.Println("*** Converted video without errors")
	}

	videoMessage := tgbotapi.NewVideo(update.Message.Chat.ID, tgbotapi.FilePath(fileName))
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
