package myerrors

import (
	"errors"
	"fmt"
)

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
	UnableToExtractWebpageVideoData    = "Unable to extract webpage video data"
)

type CookieExpiredError struct {
	Platform string
}

func (e *CookieExpiredError) Error() string {
	return RequestedContentIsNotAvailableText
}

func (e *CookieExpiredError) Severity() ErrorSeverity {
	return SeverityMaintainer
}

func (e *CookieExpiredError) UserMessage() string {
	return fmt.Sprintf("Temporary problem downloading from %s. Please try again later.", e.Platform)
}

func (e *CookieExpiredError) MaintainerMessage() string {
	return fmt.Sprintf("ALERT: %s cookies истекли. Требуется обновление cookies.",
		e.Platform)
}

type VideoUnavailableError struct {
	Platform string
}

func (e *VideoUnavailableError) Error() string {
	return VideoUnavailableText
}

func (e *VideoUnavailableError) Severity() ErrorSeverity {
	return SeverityUser
}

func (e *VideoUnavailableError) UserMessage() string {
	return "Video unavailable on " + e.Platform
}

func (e *VideoUnavailableError) MaintainerMessage() string {
	return ""
}

type UnsupportedURLError struct {
	URL      string
	Platform string
}

func (e *UnsupportedURLError) Error() string {
	return UnsupportedURLText
}

func (e *UnsupportedURLError) Severity() ErrorSeverity {
	return SeverityUser
}

func (e *UnsupportedURLError) UserMessage() string {
	return fmt.Sprintf("Unsupported URL on %s: %s", e.Platform, e.URL)
}

func (e *UnsupportedURLError) MaintainerMessage() string {
	return ""
}

type RequestEntityTooLargeError struct{}

func (e *RequestEntityTooLargeError) Error() string {
	return RequestEntityTooLargeText
}

func (e *RequestEntityTooLargeError) Severity() ErrorSeverity {
	return SeverityUser
}

func (e *RequestEntityTooLargeError) UserMessage() string {
	return "File is too large to download"
}

func (e *RequestEntityTooLargeError) MaintainerMessage() string {
	return ""
}

type UnableToExtractWebpageVideoDataError struct{}

func (e *UnableToExtractWebpageVideoDataError) Error() string {
	return UnableToExtractWebpageVideoData
}

func (e *UnableToExtractWebpageVideoDataError) Severity() ErrorSeverity {
	return SeverityUser
}

func (e *UnableToExtractWebpageVideoDataError) UserMessage() string {
	return UnableToExtractWebpageVideoData
}

func (e *UnableToExtractWebpageVideoDataError) MaintainerMessage() string {
	return ""
}
