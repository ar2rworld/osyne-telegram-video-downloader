package handler

import (
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/wader/goutubedl"

	"github.com/ar2rworld/golang-telegram-video-downloader/app/botservice"
	"github.com/ar2rworld/golang-telegram-video-downloader/app/downloader"
	"github.com/ar2rworld/golang-telegram-video-downloader/app/match"
)

type Handler struct {
	bot                  *tgbotapi.BotAPI
	botService           *botservice.BotService
	CookiesPath          string // Added this field
	InstagramCookiesPath string
	GoogleCookiesPath    string
}

func NewHandler(bot *tgbotapi.BotAPI, botService *botservice.BotService, c, i, g string) *Handler {
	return &Handler{
		bot:                  bot,
		botService:           botService,
		CookiesPath:          c,
		InstagramCookiesPath: i,
		GoogleCookiesPath:    g,
	}
}

func (h *Handler) HandleError(u *tgbotapi.Update, err error) {
	if err != nil {
		log.Println(err)
		err = h.botService.Log(u, err)
		if err != nil {
			log.Println(err)
		}
	}
}

func (h *Handler) VideoMessage(u *tgbotapi.Update, url string) error {
	remove := []string{}
	defer removeFiles(remove)

	log.Printf("*** Got request to download video: %s", url)
	opts := goutubedl.Options{HTTPClient: &http.Client{}, DebugLog: log.Default()}
	isYoutubeVideo := match.Youtube(url) != ""
	do := &goutubedl.DownloadOptions{}
	if isYoutubeVideo {
		do = alterDownloadOptions(u, url, &opts)
	}

	var fileName string
	var err error

	h.setupCookies(url, &opts, isYoutubeVideo)

	fileName, err = downloader.DownloadVideo(url, opts, do)
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

	err = h.handleAudioVideoMessage(do, u, fileName)
	if err != nil {
		return err
	}

	log.Println("*** Finished sending video/audio")
	removeFiles(remove)
	return nil
}

// If handler has instagram cookies download with cookies
// If handler has google cookies download youtube video or short with cookies
// else just try downloading
func (h *Handler) setupCookies(url string, opts *goutubedl.Options, isYoutubeVideo bool) {
	switch {
	case match.Instagram(url) != "" && h.InstagramCookiesPath != "":
		log.Println("*** Downloading Instagram with Cookies")
		opts.Cookies = h.InstagramCookiesPath
	case isYoutubeVideo && h.GoogleCookiesPath != "":
		log.Println("*** Downloading Youtube Video with Cookies")
		opts.Cookies = h.GoogleCookiesPath
	case match.YoutubeShorts(url) != "" && h.GoogleCookiesPath != "":
		log.Println("*** Downloading Youtube Shorts with Cookies")
		opts.Cookies = h.GoogleCookiesPath
	default:
		log.Println("*** DownloadVideo")
	}
}

func (h *Handler) handleAudioVideoMessage(do *goutubedl.DownloadOptions, u *tgbotapi.Update, fileName string) error {
	var err error
	if do.DownloadAudioOnly {
		log.Println("*** Started sending audio")

		audioMessage := tgbotapi.NewAudio(u.Message.Chat.ID, tgbotapi.FilePath(fileName))
		audioMessage.ReplyParameters.MessageID = u.Message.MessageID

		_, err = h.bot.Send(audioMessage)
	} else {
		log.Println("*** Started sending video")

		videoMessage := tgbotapi.NewVideo(u.Message.Chat.ID, tgbotapi.FilePath(fileName))
		videoMessage.ReplyParameters.MessageID = u.Message.MessageID

		_, err = h.bot.Send(videoMessage)
	}

	return err
}

func removeFiles(files []string) {
	for _, fn := range files {
		err := os.Remove(fn)
		log.Println("*** Removed file: ", fn, "error:", err)
	}
}
