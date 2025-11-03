package utils

import c "github.com/ar2rworld/golang-telegram-video-downloader/app/constants"

func BytesToMb(b float64) float64 {
	return b / c.BytesInKByte / c.BytesInKByte
}
