package handler

import (
	"context"
	"errors"
	"log"
	"strings"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/wader/goutubedl"

	"github.com/ar2rworld/golang-telegram-video-downloader/app/match"
	"github.com/ar2rworld/golang-telegram-video-downloader/app/myerrors"
)

func (h *Handler) HandleUpdate(ctx context.Context, wg *sync.WaitGroup, update *tgbotapi.Update) { //nolint: gocyclo, cyclop
	defer wg.Done()

	if update.Message == nil {
		return
	}

	messageText := update.Message.Text

	if len(update.Message.Entities) > 0 && update.Message.ReplyToMessage != nil && strings.Contains(messageText, h.bot.Self.UserName) {
		err := h.HandleMentionMessage(ctx, update)
		if err != nil && errors.Is(err, ErrNoURLFound) {
			err = h.Whaat(update)
			h.HandleError(update, err)
		} else if err != nil {
			h.HandleError(update, err)
		}

		return
	}

	// Inside the main loop where you handle updates
	if update.Message.From.ID == h.AdminID && update.Message.Document != nil {
		err := h.HandleAdminMessage(ctx, update)
		h.HandleError(update, err)

		return
	}

	url := match.Match(messageText)

	switch {
	case url != "":
		err := h.VideoMessage(ctx, update, url)
		if err != nil {
			h.HandleError(update, err)

			err := h.ThumbDown(update)
			if err != nil {
				log.Println("Error while reacting:", err)
			}
		}

	case messageText == "/test":
		err := h.handleAudioVideoMessage(&goutubedl.DownloadOptions{}, update, "output.mp4")
		log.Println(err.Error())
		h.HandleError(update, err)

	case messageText == "osyndaisyn ba?":
		message := tgbotapi.NewMessage(update.Message.Chat.ID, "osyndaymyn")
		sentMessage, err := h.bot.Send(message)
		myerrors.CheckTextMessage(&message, err, &sentMessage)

	case update.Message.Chat.ID == update.Message.From.ID:
		// No URL found in private message
		err := h.Whaat(update)
		h.HandleError(update, err)
	}
}
