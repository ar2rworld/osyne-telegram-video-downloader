package utilities

import (
	"strconv"

	c "github.com/ar2rworld/golang-telegram-video-downloader/app/constants"
)

func BytesToMb(b float64) float64 {
	return b / c.BytesInKByte / c.BytesInKByte
}

func ConvertSecondsToMinSec(seconds int) string {
	minutes := seconds / c.SecondsInMinute
	seconds %= c.SecondsInMinute

	return strconv.Itoa(minutes) + ":" + strconv.Itoa(seconds)
}
