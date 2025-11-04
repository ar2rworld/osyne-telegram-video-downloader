package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/ar2rworld/golang-telegram-video-downloader/app/botservice"
	"github.com/ar2rworld/golang-telegram-video-downloader/app/downloader"
	"github.com/ar2rworld/golang-telegram-video-downloader/app/handler"
	"github.com/ar2rworld/golang-telegram-video-downloader/app/logger"
	"github.com/ar2rworld/golang-telegram-video-downloader/app/platform"
)

func main() {
	botAPI, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatalln("BOT_TOKEN: ", err)
	}

	adminID, err := strconv.ParseInt(os.Getenv("ADMIN_ID"), 10, 64)
	if err != nil {
		log.Fatalln("parsing ADMIN_ID: ", err)
	}

	logChannelID, err := strconv.ParseInt(os.Getenv("LOG_CHANNEL_ID"), 10, 64)
	if err != nil {
		log.Fatalln("parsing LOG_CHANNEL_ID: ", err)
	}

	ytdlpPath := os.Getenv("YT_DLP_PATH")
	cookiesPath := os.Getenv("COOKIES_PATH")
	instagramCookiesPath := os.Getenv("INSTAGRAM_COOKIES_PATH")
	googleCookiesPath := os.Getenv("GOOGLE_COOKIES_PATH")

	registry := platform.NewRegistry()
	instagram := platform.NewInstagram(instagramCookiesPath)
	youtube := platform.NewYoutube(googleCookiesPath)
	shorts := platform.NewYoutubeShorts(googleCookiesPath)
	facebookreels := platform.NewFacebookReels()

	registry.Register(instagram)
	registry.Register(youtube)
	registry.Register(shorts)
	registry.Register(facebookreels)

	botAPI.Debug = false
	l := logger.New(os.Getenv("PROD") == "1")

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := botAPI.GetUpdatesChan(updateConfig)

	// hello message to admin
	helloMessage := tgbotapi.NewMessage(adminID, "Hello, boss")

	_, err = botAPI.Send(helloMessage)
	if err != nil {
		l.Fatal().Err(err).Msg("sending hello message")
	}

	botService := botservice.NewBotService(l, botAPI, logChannelID)

	d := downloader.NewDownloader(l, ytdlpPath)
	h := handler.NewHandler(l, botAPI, botService, registry, d, cookiesPath, instagramCookiesPath, googleCookiesPath, adminID)

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
