package handler

import (
	"fmt"
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/wader/goutubedl"

	"github.com/ar2rworld/golang-telegram-video-downloader/app/downloader"
	"github.com/ar2rworld/golang-telegram-video-downloader/app/match"
)

func setDownloadSections(opts *goutubedl.Options, start, finish int) *goutubedl.Options {
	s := downloader.ConvertSecondsToMinSec(start)
	f := downloader.ConvertSecondsToMinSec(finish)
	opts.DownloadSections = fmt.Sprintf("*%s-%s", s, f)

	return opts
}

// AlterDownloadOptions modifies the goutubedl.Options
// if arguments provided else shortens video to 30 sec
// return DownloadOptions to use goutubedl.DownloadResult.DownloadWithOptions method
func AlterDownloadOptions(u *tgbotapi.Update, url string, opts *goutubedl.Options) *goutubedl.DownloadOptions {
	// modity this one to match -s and -x
	sections := match.DownloadSectionsArgument(u.Message.Text)
	audioOnly := match.DownloadAudioArgument(u.Message.Text)

	currentTime := parseCurrentTime(url)

	if sections != "" {
		userOptions, err := parse(sections)
		if err != nil {
			log.Printf("*** Error parsing video options: %s", err.Error())

			*userOptions.Sections = downloader.DefaultSections
		}

		opts.DownloadSections = *userOptions.Sections
		log.Printf("*** Downloading section of the video: %s", opts.DownloadSections)
	} else if currentTime != "" {
		t, err := strconv.Atoi(currentTime)
		if err != nil {
			log.Printf("*** Error converting to int while changing DownloadSections for youtube: %s", err.Error())
		}

		setDownloadSections(opts, t, t+downloader.HalfMinute)
		log.Printf("*** Downloading section(%s) of the video from currentTime: %s", opts.DownloadSections, currentTime)
	}

	if opts.DownloadSections == "" && audioOnly == "" {
		opts.DownloadSections = downloader.DefaultSections
		log.Printf("*** Downloading default sections of the video: %s", downloader.DefaultSections)
	}

	do := &goutubedl.DownloadOptions{}

	if audioOnly != "" {
		log.Printf("*** Extracting audio from the video")

		do.DownloadAudioOnly = true
	}

	return do
}
