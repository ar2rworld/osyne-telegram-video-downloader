package myerrors

import "errors"

var (
	ErrPlatform           = errors.New("platform error")
	ErrDownload           = errors.New("download error")
	ErrInvalidInput       = errors.New("invalid input")
	ErrCalculatedDuration = errors.New("invalid calculated duration")
	ErrNoSuitableFormat   = errors.New("no suitable format found")
	ErrNoSizeInfo         = errors.New("no size info found")
)

// Business logic messages
const (
	RequestEntityTooLarge              = "Request Entity Too Large"
	RequestEntityTooLargeText          = "request entity too large"
	UnsupportedURL                     = "Unsupported URL"
	UnsupportedURLText                 = "unsupported URL"
	VideoUnavailable                   = "Video unavailable"
	VideoUnavailableText               = "video unavailable"
	RequestedContentIsNotAvailable     = "--cookies"
	RequestedContentIsNotAvailableText = "requested content is not available"
)

// Business logic errors
var (
	ErrRequestEntityTooLarge          = errors.New(RequestEntityTooLargeText)
	ErrUnsupportedURL                 = errors.New(UnsupportedURLText)
	ErrVideoUnavailable               = errors.New(VideoUnavailableText)
	ErrRequestedContentIsNotAvailable = errors.New(RequestedContentIsNotAvailableText)
)
