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
	"github.com/ar2rworld/golang-telegram-video-downloader/app/logger"
	"github.com/ar2rworld/golang-telegram-video-downloader/app/match"
	"github.com/ar2rworld/golang-telegram-video-downloader/app/myerrors"
	"github.com/ar2rworld/golang-telegram-video-downloader/app/platform"
)

type Handler struct {
	Logger               *logger.Logger
	bot                  *tgbotapi.BotAPI
	botService           *botservice.BotService
	CookiesPath          string // Added this field
	InstagramCookiesPath string
	GoogleCookiesPath    string
	AdminID              int64
	PlatformRegistry     *platform.Registry
	Downloader           *downloader.Downloader
}

func NewHandler(l *logger.Logger, bot *tgbotapi.BotAPI, botService *botservice.BotService, r *platform.Registry, d *downloader.Downloader, c, i, g string, adminID int64) *Handler {
	return &Handler{
		Logger:               l,
		bot:                  bot,
		botService:           botService,
		CookiesPath:          c,
		InstagramCookiesPath: i,
		GoogleCookiesPath:    g,
		AdminID:              adminID,
		PlatformRegistry:     r,
		Downloader:           d,
	}
}

func (h *Handler) HandleError(u *tgbotapi.Update, err error) {
	if err == nil {
		return
	}

	// catches json: cannot unmarshal bool into Go value of type tgbotapi.Message
	if strings.Contains(err.Error(), "cannot unmarshal bool") {
		return
	}

	// if error accured in private message, let user know that there is an error
	if u.Message != nil && u.Message.Chat.ID == u.Message.From.ID {
		var classified myerrors.ClassifiedError
		if errors.As(err, &classified) {
			msg := tgbotapi.NewMessage(u.Message.Chat.ID, classified.UserMessage())

			_, sendErr := h.bot.Send(msg)
			if sendErr != nil {
				h.Logger.Error().Err(sendErr).Msg("error sending message")
			}

			switch classified.Severity() {
			case myerrors.SeverityMaintainer:
				h.botService.AlertAdmin(classified.MaintainerMessage())
				h.Logger.Error().Msg(classified.MaintainerMessage())
			case myerrors.SeverityCritical:
				h.botService.AlertAdmin("CRITICAL: " + classified.MaintainerMessage())
				h.Logger.Error().Msg(classified.MaintainerMessage())
			case myerrors.SeverityUser:
				h.botService.Log(classified.Error())
			}
		} else {
			_, sendErr := h.bot.Send(tgbotapi.NewMessage(u.Message.Chat.ID, myerrors.InternalErrorText))
			if sendErr != nil {
				h.Logger.Error().Err(err).Msg("error sending message")
			}
		}
	}

	h.Logger.Error().Err(err).Msg("handle error update")

	h.botService.LogErrorUpdate(u, err)
}

func (h *Handler) VideoMessage(ctx context.Context, u *tgbotapi.Update, url string) error {
	prms := downloader.NewParameters()
	defer h.removeFiles(prms.TempFiles)

	h.Logger.Info().Str("url", url).Msg("Got request to download video")

	opts := goutubedl.Options{HTTPClient: &http.Client{}, DebugLog: log.Default()}
	isYoutubeVideo := match.Youtube(url) != ""

	do := &goutubedl.DownloadOptions{}
	if isYoutubeVideo {
		do = h.AlterDownloadOptions(u, url, &opts)
	}

	var (
		fileName string
		err      error
	)

	p := h.PlatformRegistry.FindPlatform(url)
	h.Logger.Info().Str("platform", p.Name()).Msg("registry found")
	prms.Platform = p

	p.ConfigureDownload(url, &opts)

	fileName, err = h.Downloader.DownloadVideo(ctx, url, opts, do, prms)
	if err != nil {
		return fmt.Errorf("error downloading video: %w", err)
	}

	prms.AddTempFile(fileName)

	h.Logger.Info().Msg("Downloaded video without errors")

	err = h.handleAudioVideoMessage(do, u, fileName)
	if err != nil {
		return fmt.Errorf("error sending video/audio: %w", err)
	}

	h.Logger.Info().Msg("Finished sending video/audio")
	h.removeFiles(prms.TempFiles)

	return nil
}

func (h *Handler) handleAudioVideoMessage(do *goutubedl.DownloadOptions, u *tgbotapi.Update, fileName string) error {
	var err error
	// TODO: upload as file document

	if do.DownloadAudioOnly {
		h.Logger.Info().Msg("Started sending audio")

		doc := tgbotapi.NewDocument(u.Message.Chat.ID, tgbotapi.FilePath(fileName))
		doc.ReplyParameters.MessageID = u.Message.MessageID

		_, err = h.bot.Send(doc)
		if err != nil {
			return fmt.Errorf("failed to send document: %w", err)
		}
	} else {
		h.Logger.Info().Msg("Started sending video")

		videoMessage := tgbotapi.NewVideo(u.Message.Chat.ID, tgbotapi.FilePath(fileName))
		videoMessage.ReplyParameters.MessageID = u.Message.MessageID

		_, err = h.bot.Send(videoMessage)
	}

	if err != nil && strings.Contains(err.Error(), myerrors.RequestEntityTooLarge) {
		return &myerrors.RequestEntityTooLargeError{}
	}

	return err
}

func (h *Handler) removeFiles(files *[]string) {
	for _, fn := range *files {
		err := os.Remove(fn)
		h.Logger.Info().Str("file", fn).Err(err).Msg("removed")
	}
}
