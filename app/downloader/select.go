package downloader

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ar2rworld/golang-telegram-video-downloader/app/constants"
	"github.com/wader/goutubedl"
)

const (
	MinHeight    = 300
	MinWidth     = 600
	BytesInKByte = 1024
)

const (
	AudioCodec    = "mp4a"
	VideoCodec    = "avc1"
)

// SelectFormat: select video format out of minimal height and width of the video
// sortout video and audio formats based on the combined filesize smaller than TgUploadLimit
// If none formats found, will return a single format containing A and V codecs and filesize
func SelectFormat(formats []goutubedl.Format) (filter, ext string, err error) { //nolint: nonamedreturns
	videoFormats := sortoutMinHeightWidth(formats, MinHeight, MinWidth)
	audioFormats := sortoutFormat(formats, AudioCodec)
	videoFormats = sortoutFormat(videoFormats, VideoCodec)

	if len(videoFormats) == 0 || len(audioFormats) == 0 {
		return "", "", ErrNoSuitableFormat
	}

	// brute force video and audio formats finding the best quality under TgUploadLimit
	bestAudioFormat := ""
	bestVideoFormat := ""
	extOut := ""

	for videoFormatIndex := range videoFormats {
		for audioFormatIndex := range audioFormats {
			audioFormat := &audioFormats[audioFormatIndex]
			videoFormat := &videoFormats[videoFormatIndex]

			currentSize := bytesToMb(videoFormat.Filesize) + bytesToMb(audioFormat.Filesize)
			if currentSize < constants.TgUploadLimit {
				bestAudioFormat = audioFormat.FormatID
				bestVideoFormat = videoFormat.FormatID
				extOut = videoFormat.Ext
			}
		}
	}

	if bestAudioFormat == "" || bestVideoFormat == "" {
		return "", "", ErrNoSuitableFormat
	}

	filterOut := fmt.Sprintf("%s+%s", bestVideoFormat, bestAudioFormat)

	return filterOut, extOut, nil
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

func bytesToMb(b float64) float64 {
	return b / BytesInKByte / BytesInKByte
}
