package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/wader/goutubedl"

	"github.com/ar2rworld/golang-telegram-video-downloader/app/botservice"
	"github.com/ar2rworld/golang-telegram-video-downloader/app/downloader"
	"github.com/ar2rworld/golang-telegram-video-downloader/app/match"
	"github.com/ar2rworld/golang-telegram-video-downloader/app/myerrors"
)

const (
	UnableToExtractWebpageVideoData = "Unable to extract webpage video data"
)

type Handler struct {
	bot                  *tgbotapi.BotAPI
	botService           *botservice.BotService
	CookiesPath          string // Added this field
	InstagramCookiesPath string
	GoogleCookiesPath    string
	AdminID              int64
}

func NewHandler(bot *tgbotapi.BotAPI, botService *botservice.BotService, c, i, g string, adminID int64) *Handler {
	return &Handler{
		bot:                  bot,
		botService:           botService,
		CookiesPath:          c,
		InstagramCookiesPath: i,
		GoogleCookiesPath:    g,
		AdminID:              adminID,
	}
}

func (h *Handler) HandleError(u *tgbotapi.Update, err error) {
	if err == nil {
		return
	}

	estr := err.Error()

	// catches json: cannot unmarshal bool into Go value of type tgbotapi.Message
	if strings.Contains(estr, "cannot unmarshal bool") {
		return
	}

	// if error accured in private message, let user know that there is an error
	if u.Message != nil && u.Message.Chat.ID == u.Message.From.ID { //nolint: nestif
		if strings.Contains(estr, UnableToExtractWebpageVideoData) {
			msg := tgbotapi.NewMessage(u.Message.Chat.ID, UnableToExtractWebpageVideoData)

			_, sendErr := h.bot.Send(msg)
			if sendErr != nil {
				log.Println(sendErr)
			}
		}

		msgText := h.selectErrorMessage(err)
		if msgText != "" {
			msg := tgbotapi.NewMessage(u.Message.Chat.ID, msgText)

			_, sendErr := h.bot.Send(msg)
			if sendErr != nil {
				log.Println(sendErr)
			}
		}

		msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Something went wrong, I will let the Creator know")

		_, sendErr := h.bot.Send(msg)
		if sendErr != nil {
			log.Println(sendErr)
		}
	}

	log.Println(err)

	h.botService.Log(u, err)
}

func (h *Handler) VideoMessage(ctx context.Context, u *tgbotapi.Update, url string) error {
	prms := downloader.NewParameters()
	defer removeFiles(prms.TempFiles)

	log.Printf("*** Got request to download video: %s", url)
	opts := goutubedl.Options{HTTPClient: &http.Client{}, DebugLog: log.Default()}
	isYoutubeVideo := match.Youtube(url) != ""

	do := &goutubedl.DownloadOptions{}
	if isYoutubeVideo {
		do = AlterDownloadOptions(u, url, &opts)
	}

	var (
		fileName string
		err      error
	)

	h.setupCookies(url, &opts, prms, isYoutubeVideo)

	fileName, err = downloader.DownloadVideo(ctx, url, opts, do, prms)
	if err != nil {
		return fmt.Errorf("error downloading video: %w", err)
	}

	prms.AddTempFile(fileName)

	log.Println("*** Downloaded video without errors")

	err = h.handleAudioVideoMessage(do, u, fileName)
	if err != nil {
		return fmt.Errorf("error sending video/audio: %w", err)
	}

	log.Println("*** Finished sending video/audio")
	removeFiles(prms.TempFiles)

	return nil
}

// If handler has instagram cookies download with cookies
// If handler has google cookies download youtube video or short with cookies
// else just try downloading
func (h *Handler) setupCookies(url string, opts *goutubedl.Options, prms *downloader.Parameters, isYoutubeVideo bool) {
	switch {
	case match.Instagram(url) != "" && h.InstagramCookiesPath != "":
		log.Println("*** Downloading Instagram with Cookies")

		prms.IsInstagram = true
		opts.Cookies = h.InstagramCookiesPath
	case isYoutubeVideo && h.GoogleCookiesPath != "":
		log.Println("*** Downloading Youtube Video with Cookies")

		prms.IsYoutubeVideo = true
		opts.Cookies = h.GoogleCookiesPath
	case match.YoutubeShorts(url) != "" && h.GoogleCookiesPath != "":
		log.Println("*** Downloading Youtube Shorts with Cookies")

		prms.IsYoutubeShorts = true
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

	if err != nil && strings.Contains(err.Error(), myerrors.RequestEntityTooLarge) {
		return myerrors.ErrRequestEntityTooLarge
	}

	return err
}

func removeFiles(files *[]string) {
	for _, fn := range *files {
		err := os.Remove(fn)
		log.Println("*** Removed file: ", fn, "error:", err)
	}
}

func (h *Handler) selectErrorMessage(err error) string {
	if errors.Is(err, myerrors.ErrRequestEntityTooLarge) {
		return "File is too large to download"
	}

	if errors.Is(err, myerrors.ErrUnsupportedURL) {
		return myerrors.UnsupportedURL
	}

	if errors.Is(err, myerrors.ErrVideoUnavailable) {
		return myerrors.VideoUnavailable
	}

	if errors.Is(err, myerrors.ErrRequestedContentIsNotAvailable) {
		return "Cookies expired, the dev will need to refresh them"
	}

	return ""
}
