package handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) ThumbDown(u *tgbotapi.Update) error {
	thumbdown := tgbotapi.ReactionType{
		Type:  "emoji",
		Emoji: "👎",
	}
	return h.reaction(u, thumbdown)
}

func (h *Handler) Whaat(u *tgbotapi.Update) error {
	whaat := tgbotapi.ReactionType{
		Type:  "emoji",
		Emoji: "🤔",
	}
	return h.reaction(u, whaat)
}

func (h *Handler) reaction(update *tgbotapi.Update, reaction tgbotapi.ReactionType) error {
	r := tgbotapi.SetMessageReactionConfig{
		BaseChatMessage: tgbotapi.BaseChatMessage{
			ChatConfig: tgbotapi.ChatConfig{
				ChatID: update.Message.Chat.ID,
			},
			MessageID: update.Message.MessageID,
		},
		Reaction: []tgbotapi.ReactionType{
			reaction,
		},
	}

	_, err := h.bot.Send(r)

	return err
}
