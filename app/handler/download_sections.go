package handler

import (
	"fmt"
	"log"
	"strconv"

	"github.com/ar2rworld/golang-telegram-video-downloader/app/match"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/wader/goutubedl"
)

func convertSecondsToMinSec(seconds int) string {
	minutes := seconds / 60
	seconds = seconds % 60
	return strconv.Itoa(minutes) + ":" + strconv.Itoa(seconds)
}

func setDownloadSections(opts *goutubedl.Options, start, finish int) *goutubedl.Options {
	s := convertSecondsToMinSec(start)
	f := convertSecondsToMinSec(finish)
	opts.DownloadSections = fmt.Sprintf("*%s-%s", s, f)
	return opts
}


func alterDownloadSections(u *tgbotapi.Update, url string, opts *goutubedl.Options) {
	args := match.DownloadSectionsArgument(u.Message.Text)
	currentTime := parseCurrentTime(url)

	if args != "" {
		sections, err := parse(args)
		if err != nil {
			log.Printf("*** Error parsing video sections: %s", err.Error())
			sections = DefaultSections
		}
		opts.DownloadSections = sections
		log.Printf("*** Downloading section of the video: %s", opts.DownloadSections)
	} else if currentTime != "" {
		t, err := strconv.Atoi(currentTime)
		if err != nil {
			log.Printf("*** Error converting to int while changing DownloadSections for youtube: %s", err.Error())
		}
		setDownloadSections(opts, t, t + 30)
		log.Printf("*** Downloading section(%s) of the video from currentTime: %s", opts.DownloadSections, currentTime)
	}

	if opts.DownloadSections == "" {
		opts.DownloadSections = DefaultSections
		log.Printf("*** Downloading default sections of the video: %s", DefaultSections)
	}
}
