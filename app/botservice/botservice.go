package botservice

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotService struct {
	api          *tgbotapi.BotAPI
	logChannelID int64
}

func NewBotService(api *tgbotapi.BotAPI, logChannelID int64) *BotService {
	return &BotService{
		api:          api,
		logChannelID: logChannelID,
	}
}

func (b *BotService) Log(u *tgbotapi.Update, err error) error {
	if u != nil && u.Message == nil {
		log.Println("BotService Log used without message in update")
		return nil
	}

	if err == nil {
		log.Println("BotService Log used without error")
		return nil
	}

	text := "Error in " + u.Message.Chat.Title + " (" + u.Message.Chat.UserName + "): " + err.Error()

	msg := tgbotapi.NewMessage(b.logChannelID, text)
	_, err = b.api.Send(msg)

	return err
}
