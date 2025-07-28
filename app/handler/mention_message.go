package handler

import (
	"context"
	"errors"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/ar2rworld/golang-telegram-video-downloader/app/match"
)

var ErrNoURLFound = errors.New("no url found in message")

func (h *Handler) HandleMentionMessage(ctx context.Context, u *tgbotapi.Update) error {
	messageText := u.Message.Text
	if !strings.Contains(messageText, h.bot.Self.UserName) {
		return nil
	}

	log.Println("*** Mentioned Username in message")

	url := match.Match(u.Message.ReplyToMessage.Text)

	if url == "" {
		return ErrNoURLFound
	}

	return h.VideoMessage(ctx, u, url)
}
