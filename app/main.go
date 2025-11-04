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
	"github.com/ar2rworld/golang-telegram-video-downloader/app/handler"
	"github.com/ar2rworld/golang-telegram-video-downloader/app/myerrors"
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

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := botAPI.GetUpdatesChan(updateConfig)

	// hello message to admin
	helloMessage := tgbotapi.NewMessage(adminID, "Hello, boss")
	sentMessage, err := botAPI.Send(helloMessage)
	myerrors.CheckTextMessage(&helloMessage, err, &sentMessage)

	botService := botservice.NewBotService(botAPI, logChannelID)
	h := handler.NewHandler(botAPI, botService, registry, cookiesPath, instagramCookiesPath, googleCookiesPath, adminID)

	// Create a context to handle graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())

	// Create a WaitGroup to keep track of running goroutines
	wg := &sync.WaitGroup{}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Starting handling updates")

	go func() {
		for update := range updates {
			wg.Add(1)

			go h.HandleUpdate(ctx, wg, &update)
		}
	}()

	signalReceived := <-signalCh
	log.Println("signalReceived: ", signalReceived)

	log.Println("HandleUpdate loop has stopped")

	cancel()

	wg.Wait()

	log.Println("Shutdown complete.")
}
