package handler

import (
	"fmt"
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/wader/goutubedl"

	"github.com/ar2rworld/golang-telegram-video-downloader/app/match"
)

const (
	halfMinute      = 30
	secondsInMinute = 60
)

func convertSecondsToMinSec(seconds int) string {
	minutes := seconds / secondsInMinute
	seconds %= secondsInMinute
	return strconv.Itoa(minutes) + ":" + strconv.Itoa(seconds)
}

func setDownloadSections(opts *goutubedl.Options, start, finish int) *goutubedl.Options {
	s := convertSecondsToMinSec(start)
	f := convertSecondsToMinSec(finish)
	opts.DownloadSections = fmt.Sprintf("*%s-%s", s, f)
	return opts
}

// alterDownloadOptions modifies the goutubedl.Options
// if arguments provided else shortens video to 30 sec
// return DownloadOptions to use goutubedl.DownloadResult.DownloadWithOptions method
func alterDownloadOptions(u *tgbotapi.Update, url string, opts *goutubedl.Options) *goutubedl.DownloadOptions {
	// modity this one to match -s and -x
	sections := match.DownloadSectionsArgument(u.Message.Text)
	audioOnly := match.DownloadAudioArgument(u.Message.Text)

	currentTime := parseCurrentTime(url)

	if sections != "" {
		userOptions, err := parse(sections)
		if err != nil {
			log.Printf("*** Error parsing video options: %s", err.Error())
			*userOptions.Sections = DefaultSections
		}
		opts.DownloadSections = *userOptions.Sections
		log.Printf("*** Downloading section of the video: %s", opts.DownloadSections)
	} else if currentTime != "" {
		t, err := strconv.Atoi(currentTime)
		if err != nil {
			log.Printf("*** Error converting to int while changing DownloadSections for youtube: %s", err.Error())
		}
		setDownloadSections(opts, t, t+halfMinute)
		log.Printf("*** Downloading section(%s) of the video from currentTime: %s", opts.DownloadSections, currentTime)
	}

	if opts.DownloadSections == "" && audioOnly == "" {
		opts.DownloadSections = DefaultSections
		log.Printf("*** Downloading default sections of the video: %s", DefaultSections)
	}

	do := &goutubedl.DownloadOptions{}
	if audioOnly != "" {
		log.Printf("*** Extracting audio from the video")
		do.DownloadAudioOnly = true
	}

	return do
}
