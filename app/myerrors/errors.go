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

type ErrCookieExpired struct {
	Platform string
}

func (e *ErrCookieExpired) Error() string {
	return RequestedContentIsNotAvailableText
}
func (e *ErrCookieExpired) Severity() ErrorSeverity {
	return SeverityMaintainer
}
func (e *ErrCookieExpired) UserMessage() string {
	return fmt.Sprintf("Temporary problem downloading from %s. Please try again later.", e.Platform)
}
func (e *ErrCookieExpired) MaintainerMessage() string {
	return fmt.Sprintf("ALERT: %s cookies истекли. Требуется обновление cookies.",
		e.Platform)
}

type ErrVideoUnavailable struct {
	Platform string
}

func (e *ErrVideoUnavailable) Error() string {
	return VideoUnavailableText
}
func (e *ErrVideoUnavailable) Severity() ErrorSeverity {
	return SeverityUser
}
func (e *ErrVideoUnavailable) UserMessage() string {
	return fmt.Sprintf("Video unavailable on %s", e.Platform)
}
func (e *ErrVideoUnavailable) MaintainerMessage() string {
	return ""
}

type ErrUnsupportedURL struct {
	URL      string
	Platform string
}

func (e *ErrUnsupportedURL) Error() string {
	return UnsupportedURLText
}
func (e *ErrUnsupportedURL) Severity() ErrorSeverity {
	return SeverityUser
}
func (e *ErrUnsupportedURL) UserMessage() string {
	return fmt.Sprintf("Unsupported URL on %s: %s", e.Platform, e.URL)
}
func (e *ErrUnsupportedURL) MaintainerMessage() string {
	return ""
}

type ErrRequestEntityTooLarge struct{}

func (e *ErrRequestEntityTooLarge) Error() string {
	return RequestEntityTooLargeText
}
func (e *ErrRequestEntityTooLarge) Severity() ErrorSeverity {
	return SeverityUser
}
func (e *ErrRequestEntityTooLarge) UserMessage() string {
	return "File is too large to download"
}
func (e *ErrRequestEntityTooLarge) MaintainerMessage() string {
	return ""
}

type ErrUnableToExtractWebpageVideoData struct{}

func (e *ErrUnableToExtractWebpageVideoData) Error() string {
	return UnableToExtractWebpageVideoData
}
func (e *ErrUnableToExtractWebpageVideoData) Severity() ErrorSeverity {
	return SeverityUser
}
func (e *ErrUnableToExtractWebpageVideoData) UserMessage() string {
	return UnableToExtractWebpageVideoData
}
func (e *ErrUnableToExtractWebpageVideoData) MaintainerMessage() string {
	return ""
}
