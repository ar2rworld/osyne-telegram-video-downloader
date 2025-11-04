package platform

import (
	"fmt"
	"math"
	"strings"

	"github.com/wader/goutubedl"

	c "github.com/ar2rworld/golang-telegram-video-downloader/app/constants"
	"github.com/ar2rworld/golang-telegram-video-downloader/app/match"
	"github.com/ar2rworld/golang-telegram-video-downloader/app/myerrors"
	"github.com/ar2rworld/golang-telegram-video-downloader/app/utils"
)

const FileSizeFix = 0.5

type YouTube struct {
	cookiesPath string
}

// RemuxRequired implements Platform.
func (y *YouTube) RemuxRequired() bool {
	return true
}

// ConfigureDownload implements Platform.
func (y *YouTube) ConfigureDownload(url string, opts *goutubedl.Options) {
	opts.Cookies = y.cookiesPath
}

// Match implements Platform.
func (y *YouTube) Match(url string) bool {
	return match.Youtube(url) != ""
}

// Name implements Platform.
func (y *YouTube) Name() string {
	return "youtube"
}

func NewYoutube(cookiesPath string) *YouTube {
	return &YouTube{cookiesPath: cookiesPath}
}

func (y *YouTube) SelectFormat(formats []goutubedl.Format) (format string, err error) {
	videoFormats := sortoutMinHeightWidth(formats, c.MinHeight, c.MinWidth)
	completeFormats := sortoutCompleteFormat(videoFormats)

	bestFormat := selectBestFormat(completeFormats)
	if bestFormat != nil {
		return bestFormat.FormatID, nil
	}

	audioFormats := sortoutFormat(formats, c.AudioCodec)
	videoFormats = sortoutFormat(videoFormats, c.VideoCodec)

	if len(videoFormats) == 0 || len(audioFormats) == 0 {
		return "", myerrors.ErrNoSuitableFormat
	}

	// brute force video and audio formats finding the best quality under TgUploadLimit
	bestAudioFormat := ""
	bestVideoFormat := ""

	for videoFormatIndex := range videoFormats {
		for audioFormatIndex := range audioFormats {
			audioFormat := &audioFormats[audioFormatIndex]
			videoFormat := &videoFormats[videoFormatIndex]

			currentSize := utils.BytesToMb(videoFormat.Filesize) + utils.BytesToMb(audioFormat.Filesize)
			if currentSize < c.TgUploadLimit {
				bestAudioFormat = audioFormat.FormatID
				bestVideoFormat = videoFormat.FormatID
			}
		}
	}

	if bestAudioFormat == "" || bestVideoFormat == "" {
		return "", myerrors.ErrNoSuitableFormat
	}

	filterOut := fmt.Sprintf("%s+%s", bestVideoFormat, bestAudioFormat)

	return filterOut, nil
}

// MaxDuration calculates how many seconds of the video will fit under TgUploadLimit.
// args: filesize (float64 in MB), duration (float64 in seconds)
// returns: int (seconds), error if calculation invalid
func (y *YouTube) MaxDuration(r *goutubedl.Result) (string, error) {
	var filesize float64

	duration := r.Info.Duration

	if r.Info.Filesize != 0.0 {
		filesize = utils.BytesToMb(r.Info.Filesize)
	} else if r.Info.FilesizeApprox != 0.0 {
		filesize = utils.BytesToMb(r.Info.FilesizeApprox)
	} else {
		return "", fmt.Errorf("%w: filesize=%.2f MB, duration=%.2f s", myerrors.ErrInvalidInput, filesize, duration)
	}

	// Calculate ratio
	ratio := c.TgUploadLimit * FileSizeFix / filesize
	// if ratio >= 1.0 {
	// 	// Whole video fits
	// 	return nil
	// }
	// TODO: test if short video is ok

	// Calculate maximum allowed seconds
	maxSeconds := duration * ratio
	if maxSeconds < 1 {
		return "", fmt.Errorf("%w too short: %.2f s", myerrors.ErrCalculatedDuration, maxSeconds)
	}

	seconds := int(math.Floor(maxSeconds))

	return "*0:0-" + utils.ConvertSecondsToMinSec(seconds), nil
}

func sortoutMinHeightWidth(formats []goutubedl.Format, minHeight, minWidth float64) []goutubedl.Format {
	out := make([]goutubedl.Format, 0, len(formats)) // capacity optimization

	for i := range formats {
		format := &formats[i]
		if format.Height > minHeight && format.Width > minWidth {
			out = append(out, *format)
		}
	}

	return out
}

func sortoutFormat(formats []goutubedl.Format, codec string) []goutubedl.Format {
	out := make([]goutubedl.Format, 0, len(formats)) // capacity optimization

	for formatIndex := range formats {
		format := &formats[formatIndex]
		if _, hasVideoCodec := strings.CutPrefix(format.VCodec, codec); hasVideoCodec && format.ACodec == "none" {
			out = append(out, *format)
		}

		if _, hasAudioCodec := strings.CutPrefix(format.ACodec, codec); hasAudioCodec && format.VCodec == "none" {
			out = append(out, *format)
		}
	}

	return out
}

func sortoutCompleteFormat(formats []goutubedl.Format) []goutubedl.Format {
	out := make([]goutubedl.Format, 0, len(formats))
	for formatIndex, vFormat := range formats {
		format := &formats[formatIndex]
		_, hasVideoCodec := strings.CutPrefix(vFormat.VCodec, c.VideoCodec)

		_, hasAudioCodec := strings.CutPrefix(format.ACodec, c.AudioCodec)
		if hasVideoCodec && hasAudioCodec {
			out = append(out, *format)
		}
	}

	return out
}

func selectBestFormat(formats []goutubedl.Format) *goutubedl.Format {
	var bestFormat *goutubedl.Format

	for formatIndex := range formats {
		format := &formats[formatIndex]
		if format.Filesize < c.TgUploadLimit || format.FilesizeApprox < c.TgUploadLimit {
			bestFormat = format
		}
	}

	return bestFormat
}

// If filesize(Mb) is over the Telegram Bot API limit
func (y *YouTube) NeedCut(r *goutubedl.Result) (bool, error) {
	if r.Info.Filesize == 0 || r.Info.FilesizeApprox == 0 {
		return false, myerrors.ErrNoSizeInfo
	}

	return r.Info.Filesize >= c.TgUploadLimit || r.Info.FilesizeApprox >= c.TgUploadLimit, nil
}
