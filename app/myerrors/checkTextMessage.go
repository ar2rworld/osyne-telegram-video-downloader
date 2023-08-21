package myerrors

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func CheckTextMessage(messageConfig *tgbotapi.MessageConfig, err error, sentMessage *tgbotapi.Message) {
	if err != nil {
		log.Printf("Error sending message: %s with in chatId:%d\nerror: %s\nsentMessage:%s", messageConfig.Text, sentMessage.Chat.ID, err, sentMessage.Text)
	}
}
