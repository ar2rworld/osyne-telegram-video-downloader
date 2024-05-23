package handler

import (
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/wader/goutubedl"

	"github.com/ar2rworld/golang-telegram-video-downloader/app/downloader"
	"github.com/ar2rworld/golang-telegram-video-downloader/app/match"
)

func VideoMessage(update tgbotapi.Update, url, cookiesPath string, bot *tgbotapi.BotAPI) error {
	remove := []string{}
	defer func() {
		for _, fn := range remove {
			err := os.Remove(fn)
			log.Println("*** Removed file: ", fn, "error:", err)
		}
	}()

	log.Printf("*** Got request to download video: %s", url)

	opts := goutubedl.Options{HTTPClient: &http.Client{}, DebugLog: log.Default()}
	isYoutubeVideo := match.Youtube(url) != ""
	if isYoutubeVideo {
		args := match.DownloadSectionsArgument(update.Message.Text)
		sections, err := parse(args)
		if err != nil {
			log.Println("*** Error parsing video sections")
			sections = DefaultSections
		}
		opts.DownloadSections = sections
		log.Printf("*** Downloading video from Youtube %s", opts.DownloadSections)
	}

	var fileName string
	var err error

	// if Instagram and cookiesPath is defined download with cookies
	if match.Instagram(url) != "" && cookiesPath != "" {
		fileName, err = downloader.DownloadWithCookies(url, cookiesPath)
	} else {
		fileName, err = downloader.DownloadVideo(url, opts)
	}
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
