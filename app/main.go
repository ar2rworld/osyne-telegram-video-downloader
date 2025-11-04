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
	"github.com/ar2rworld/golang-telegram-video-downloader/app/downloader"
	"github.com/ar2rworld/golang-telegram-video-downloader/app/handler"
	"github.com/ar2rworld/golang-telegram-video-downloader/app/logger"
	"github.com/ar2rworld/golang-telegram-video-downloader/app/platform"
)

type Options struct {
	Prod                 bool   `description:"Prod env"                env:"PROD"                   long:"prod"                   short:"p"`
	YtdlpPath            string `description:"YT-DLP path"             env:"YT_DLP_PATH"            long:"ytdlp_path"             required:"true" short:"d"`
	BotToken             string `description:"Telegram bot token"      env:"BOT_TOKEN"              long:"bot_token"              required:"true" short:"t"`
	AdminID              int64  `description:"Telegram admin id"       env:"ADMIN_ID"               long:"admin_id"               required:"true" short:"a"`
	LogChannelID         int64  `description:"Telegram log channel id" env:"LOG_CHANNEL_ID"         long:"log_channel_id"         required:"true" short:"l"`
	CookiesPath          string `description:"Cookies path"            env:"COOKIES_PATH"           long:"cookies_path"           required:"true" short:"c"`
	InstagramCookiesPath string `description:"Instagram cookies path"  env:"INSTAGRAM_COOKIES_PATH" long:"instagram_cookies_path" required:"true" short:"i"`
	YouTubeCookiesPath   string `description:"YouTube cookies path"    env:"GOOGLE_COOKIES_PATH"    long:"youtube_cookies_path"   required:"true" short:"y"`
}

func main() {
	l := logger.New(os.Getenv("PROD") == "1")

	options := &Options{}

	_, err := flags.Parse(options)
	if err != nil {
		l.Logger.Fatal().Err(err).Msg("parsing flags")
	}

	botAPI, err := tgbotapi.NewBotAPI(options.BotToken)
	if err != nil {
		l.Logger.Fatal().Err(err).Msg("creating bot api")
	}

	registry := platform.NewRegistry()
	instagram := platform.NewInstagram(options.InstagramCookiesPath)
	youtube := platform.NewYoutube(options.YouTubeCookiesPath)
	shorts := platform.NewYoutubeShorts(options.YouTubeCookiesPath)
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
	helloMessage := tgbotapi.NewMessage(options.AdminID, "Hello, boss")

	_, err = botAPI.Send(helloMessage)
	if err != nil {
		l.Fatal().Err(err).Msg("sending hello message")
	}

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
