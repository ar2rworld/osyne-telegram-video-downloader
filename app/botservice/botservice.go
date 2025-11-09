package botservice

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/ar2rworld/golang-telegram-video-downloader/app/logger"
)

type BotService struct {
	logger       *logger.Logger
	api          *tgbotapi.BotAPI
	logChannelID int64
}

func NewBotService(l *logger.Logger, api *tgbotapi.BotAPI, logChannelID int64) *BotService {
	return &BotService{
		logger:       l,
		api:          api,
		logChannelID: logChannelID,
	}
}

func (b *BotService) Log(text string) {
	b.SendLogChannel(text, false)
}

func (b *BotService) SendLogChannel(text string, notify bool) {
	msg := tgbotapi.NewMessage(b.logChannelID, text)
	msg.DisableNotification = ! notify

	_, err := b.api.Send(msg)
	if err != nil {
		b.logger.Error().Str("text", text).Msg("BotService Log error: " + err.Error())
	}
}

func (b *BotService) LogErrorUpdate(u *tgbotapi.Update, err error) {
	if u != nil && u.Message == nil {
		b.logger.Warn().Msg("BotService Log used without message in update")
		return
	}

	if err == nil {
		b.logger.Warn().Msg("BotService Log used without error")
		return
	}

	text := fmt.Sprintf("Error in %s (%s): %s", u.Message.Chat.Title, u.Message.Chat.UserName, err.Error())
	b.Log(text)

	text = "Error msg text: " + u.Message.Text
	b.Log(text)
}
