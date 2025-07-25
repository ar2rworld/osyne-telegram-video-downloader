package handler

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"path"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) HandleAdminMessage(u *tgbotapi.Update) error {
	if u.Message.Document == nil {
		message := tgbotapi.NewMessage(u.Message.Chat.ID, "Please attach a cookies file.")
		_, err := h.bot.Send(message)

		return err
	}

	// Get file info from Telegram
	file, err := h.bot.GetFile(tgbotapi.FileConfig{FileID: u.Message.Document.FileID})
	if err != nil {
		return err
	}

	documentName := u.Message.Document.FileName
	filePath := path.Join(h.CookiesPath, documentName)

	// Create a new file
	newFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer newFile.Close()

	ctx := context.Background()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, file.Link(h.bot.Token), http.NoBody)
	if err != nil {
		return err
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Copy the downloaded content to the new file
	_, err = io.Copy(newFile, resp.Body)
	if err != nil {
		return err
	}

	log.Printf("Cookies file replaced at: %s", h.CookiesPath)

	// Send confirmation message
	message := tgbotapi.NewMessage(u.Message.Chat.ID, "Cookies file updated successfully.")
	_, err = h.bot.Send(message)

	return err
}
