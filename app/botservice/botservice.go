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

func (b *BotService) Log(u *tgbotapi.Update, err error) {
	if u != nil && u.Message == nil {
		log.Println("BotService Log used without message in update")
	}

	if err == nil {
		log.Println("BotService Log used without error")
	}

	text := "Error in " + u.Message.Chat.Title + " (" + u.Message.Chat.UserName + "): " + err.Error()

	msg := tgbotapi.NewMessage(b.logChannelID, text)

	_, err = b.api.Send(msg)
	if err != nil {
		log.Println("BotService Log error: " + err.Error())
	}

	msg = tgbotapi.NewMessage(b.logChannelID, "Error msg text: "+u.Message.Text)

	_, err = b.api.Send(msg)
	if err != nil {
		log.Println("BotService Log error: " + err.Error())
	}
}
