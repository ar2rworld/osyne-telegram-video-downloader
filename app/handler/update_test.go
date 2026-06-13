package handler

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"

	"github.com/ar2rworld/golang-telegram-video-downloader/app/logger"
)

func TestIsAdminPrivateDocument(t *testing.T) {
	const adminID int64 = 42

	tests := []struct {
		name    string
		message *tgbotapi.Message
		want    bool
	}{
		{
			name: "admin document in private chat",
			message: &tgbotapi.Message{
				From:     &tgbotapi.User{ID: adminID},
				Chat:     tgbotapi.Chat{ID: adminID},
				Document: &tgbotapi.Document{FileID: "file-id"},
			},
			want: true,
		},
		{
			name: "admin document in public chat",
			message: &tgbotapi.Message{
				From:     &tgbotapi.User{ID: adminID},
				Chat:     tgbotapi.Chat{ID: -100123},
				Document: &tgbotapi.Document{FileID: "file-id"},
			},
			want: false,
		},
		{
			name: "non-admin document in private chat",
			message: &tgbotapi.Message{
				From:     &tgbotapi.User{ID: 7},
				Chat:     tgbotapi.Chat{ID: 7},
				Document: &tgbotapi.Document{FileID: "file-id"},
			},
			want: false,
		},
		{
			name: "admin private chat without document",
			message: &tgbotapi.Message{
				From: &tgbotapi.User{ID: adminID},
				Chat: tgbotapi.Chat{ID: adminID},
			},
			want: false,
		},
		{
			name: "missing sender",
			message: &tgbotapi.Message{
				Chat:     tgbotapi.Chat{ID: adminID},
				Document: &tgbotapi.Document{FileID: "file-id"},
			},
			want: false,
		},
		{
			name:    "missing message",
			message: nil,
			want:    false,
		},
	}

	h := NewHandler(logger.New(false), nil, nil, nil, nil, "c", "i", "g", adminID)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, h.isAdminPrivateDocument(tt.message))
		})
	}
}
