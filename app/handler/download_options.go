package handler

import (
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/wader/goutubedl"

	c "github.com/ar2rworld/golang-telegram-video-downloader/app/constants"
	"github.com/ar2rworld/golang-telegram-video-downloader/app/match"
	"github.com/ar2rworld/golang-telegram-video-downloader/app/utils"
)

func setDownloadSections(opts *goutubedl.Options, start, finish int) *goutubedl.Options {
	s := utils.ConvertSecondsToMinSec(start)
	f := utils.ConvertSecondsToMinSec(finish)
	opts.DownloadSections = fmt.Sprintf("*%s-%s", s, f)

	return opts
}

// AlterDownloadOptions modifies the goutubedl.Options
// if arguments provided else shortens video to 30 sec
// return DownloadOptions to use goutubedl.DownloadResult.DownloadWithOptions method
func (h *Handler) AlterDownloadOptions(u *tgbotapi.Update, url string, opts *goutubedl.Options) *goutubedl.DownloadOptions {
	// modity this one to match -s and -x
	sections := match.DownloadSectionsArgument(u.Message.Text)
	audioOnly := match.DownloadAudioArgument(u.Message.Text)

	currentTime := parseCurrentTime(url)

	if sections != "" {
		userOptions, err := parse(sections)
		if err != nil {
			h.Logger.Error().Err(err).Msg("while parsing video options")

			*userOptions.Sections = c.DefaultSections
		}

		opts.DownloadSections = *userOptions.Sections
		h.Logger.Info().Str("downloadsections", opts.DownloadSections).Msg("downloading section of the video")
	} else if currentTime != "" {
		t, err := strconv.Atoi(currentTime)
		if err != nil {
			h.Logger.Error().Err(err).Msg("converting to int while changing DownloadSections for youtube")
		}

		setDownloadSections(opts, t, t+c.HalfMinute)
		h.Logger.Info().
			Str("downloadsections", opts.DownloadSections).
			Str("currenttime", currentTime).
			Msg("downloading section of the video from currentTime")
	}

	if opts.DownloadSections == "" && audioOnly == "" {
		opts.DownloadSections = c.DefaultSections
		h.Logger.Info().Str("downloadsections", opts.DownloadSections).Msg("downloading default sections of the video")
	}

	do := &goutubedl.DownloadOptions{}

	if audioOnly != "" {
		h.Logger.Info().Msg("extracting audio from the video")

		do.DownloadAudioOnly = true
	}

	return do
}
