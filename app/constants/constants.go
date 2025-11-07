package constants

// Telegram Bot API video upload limit of 50 Mb
const TgUploadLimit = 50

// MinHeight, MinWidth - minimal video sizes for format selection
// BytesInKByte - number of bytes in kb for video file size calculation
const (
	MinHeight    = 300
	MinWidth     = 600
	BytesInKByte = 1024
)

// yt-dlp "best" filter
const BestFormat = "best"

// Video codecs for format selection
const (
	AudioCodec = "mp4a"
	VideoCodec = "avc1"
)

// Integers used to calculate video durations
// HalfMinute - 30 seconds in min
// SecondsInMinute - 60 seconds in min
const (
	HalfMinute      = 30
	SecondsInMinute = 60
)

// MaxFileNameLength - max file name for Telegram Bot API
const (
	MaxFileNameLength = 90
	DefaultSections   = "*0:0-0:30"
)
