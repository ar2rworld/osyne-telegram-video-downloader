package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jessevdk/go-flags"

	"github.com/ar2rworld/golang-telegram-video-downloader/app/botservice"
	c "github.com/ar2rworld/golang-telegram-video-downloader/app/constants"
	"github.com/ar2rworld/golang-telegram-video-downloader/app/downloader"
	"github.com/ar2rworld/golang-telegram-video-downloader/app/handler"
	"github.com/ar2rworld/golang-telegram-video-downloader/app/logger"
	"github.com/ar2rworld/golang-telegram-video-downloader/app/platform"
)

func main() {
	l := logger.New(os.Getenv("PROD") == "1")

	options := &c.Options{}

	_, err := flags.Parse(options)
	if err != nil {
		l.Logger.Fatal().Err(err).Msg("parsing flags")
	}

	botAPI, err := tgbotapi.NewBotAPI(options.BotToken)
	if err != nil {
		l.Logger.Fatal().Err(err).Msg("creating bot api")
	}

	registry := platform.NewRegistry(options)

	botAPI.Debug = false

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := botAPI.GetUpdatesChan(updateConfig)

	// hello message to admin
	sendHelloMessage(options, botAPI, l)

	botService := botservice.NewBotService(l, botAPI, options.LogChannelID)

	d := downloader.NewDownloader(l, options.YtdlpPath)
	h := handler.NewHandler(l, botAPI, botService, registry, d,
		options.CookiesPath, options.InstagramCookiesPath, options.YouTubeCookiesPath, options.AdminID)

	// Create a context to handle graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())

	// Create a WaitGroup to keep track of running goroutines
	wg := &sync.WaitGroup{}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	l.Info().Msg("Starting handling updates")

	go func() {
		for update := range updates {
			wg.Add(1)

			go h.HandleUpdate(ctx, wg, &update)
		}
	}()

	signalReceived := <-signalCh
	l.Info().Any("signal", signalReceived).Msg("signal received")

	l.Info().Msg("HandleUpdate loop has stopped")

	cancel()

	wg.Wait()

	l.Info().Msg("Shutdown complete.")
}

func sendHelloMessage(options *c.Options, botAPI *tgbotapi.BotAPI, l *logger.Logger) {
	helloMessage := tgbotapi.NewMessage(options.AdminID, "Hello, boss")

	_, err := botAPI.Send(helloMessage)
	if err != nil {
		l.Fatal().Err(err).Msg("sending hello message")
	}
}
