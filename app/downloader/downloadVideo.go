package downloader

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/wader/goutubedl"

	"github.com/ar2rworld/golang-telegram-video-downloader/app/myerrors"
)

const (
	MaxFileNameLength = 90
	DefaultSections   = "*0:0-0:30"
)

const (
	HalfMinute      = 30
	SecondsInMinute = 60
)
const FileSizeFix = 0.6

type Parameters struct {
	IsYoutubeVideo  bool
	IsInstagram     bool
	IsYoutubeShorts bool
	TempFiles       *[]string
}

func NewParameters() *Parameters {
	return &Parameters{
		TempFiles: &[]string{},
	}
}

func (p *Parameters) AddTempFile(s string) {
	*p.TempFiles = append(*p.TempFiles, s)
}

// Downloads video with options and writes to file
// If the video is a youtube video, then try to find a format (audio + video codecs) that fit in TgUploadLimit
// If not, download best quality of the video downloadSection calculated with MaxDuration fitting in TgUploadLimit * FileSizeFix
// After download If youtube video, remux video to MP4
// If video ext is not mp4, Convert file
func DownloadVideo(ctx context.Context, url string, opts goutubedl.Options, do *goutubedl.DownloadOptions, prms *Parameters) (string, error) { //nolint: gocyclo,cyclop,funlen
	goutubedl.Path = "yt-dlp_macos"

	isDefaultSection := opts.DownloadSections == DefaultSections
	// if DefaultSections is set, select video section under TgUploadLimit
	if isDefaultSection {
		opts.DownloadSections = ""
	}

	result, err := goutubedl.New(ctx, url, opts)
	if err != nil {
		return "", err
	}

	if do == nil {
		do = &goutubedl.DownloadOptions{}
	}

	isYoutubeVideo := prms.IsYoutubeVideo
	ext := ""
	cutRequired := false

	if isYoutubeVideo {
		filter, e, err := SelectFormat(result.Formats())
		if err != nil && !errors.Is(err, ErrNoSuitableFormat) {
			return "", err
		}

		if err != nil {
			do.Filter = "best"
			cutRequired = true
		}

		// if can't find suitable format , cut the video accordingly to TgUploadLimit limit
		ext = e
		do.Filter = filter
		log.Println("*** filter:", filter)
	}

	if cutRequired && isDefaultSection {
		seconds, err := MaxDuration(bytesToMb(result.Info.FilesizeApprox), result.Info.Duration)
		if err != nil {
			return "", fmt.Errorf("failed to figure out maxduration of video: %w", err)
		}

		opts.DownloadSections = "*0:0-" + ConvertSecondsToMinSec(seconds)

		result, err = goutubedl.New(ctx, url, opts)
		if err != nil {
			return "", err
		}
	}

	// Weird stuff to make yt-dlp download only audio
	if !do.DownloadAudioOnly && do.Filter == "" {
		do.Filter = "best"
	}

	downloadResult, err := result.DownloadWithOptions(ctx, *do)
	if err != nil {
		return "", err
	}
	defer downloadResult.Close()

	title := RemoveNonAlphanumericRegex(ConvertToUTF8(result.Info.Title))
	output := cutString(title, MaxFileNameLength) + result.Info.Ext

	f, err := os.CreateTemp("", output)
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = io.Copy(f, downloadResult)
	if err != nil {
		return "", err
	}

	filename := f.Name()
	if isYoutubeVideo {
		filename, err = RemuxToMP4(ctx, filename)
		if err != nil {
			return "", fmt.Errorf("could not remux video to mp4: %w", err)
		}

		prms.AddTempFile(filename)

		if ext != "mp4" || result.Info.Ext != "mp4" {
			filename, err = Convert(ctx, filename)
			if err != nil {
				return "", err
			}

			prms.AddTempFile(filename)

			log.Println("*** Converted video without errors")
		}
	}

	return filename, nil
}

func DownloadWithCookies(ctx context.Context, url, cookiesPath string) (string, error) {
	fileName := "videoDownloadedWithCookies"
	cmd := exec.CommandContext(ctx, "yt-dlp", "-f", "mp4", "-o", fileName, "--cookies", cookiesPath, url)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return fileName, nil
}

func cutString(s string, maxLength int) string {
	absMaxLength := int(math.Abs(float64(maxLength)))
	if len(s) <= absMaxLength {
		return s
	}

	return s[:absMaxLength]
}

func ConvertToUTF8(s string) string {
	result := make([]byte, 0, len(s))
	for s != "" {
		r, size := utf8.DecodeRuneInString(s)
		if r != utf8.RuneError || size > 1 {
			result = append(result, s[:size]...)
		}

		s = s[size:]
	}

	return string(result)
}

// RemoveNonAlphanumericRegex uses a regular expression to remove non-alphanumeric characters.
// This version also preserves spaces but trims leading and trailing whitespace.
func RemoveNonAlphanumericRegex(s string) string {
	reg := regexp.MustCompile(`[^a-zA-Z0-9\s]+`)
	return strings.TrimSpace(reg.ReplaceAllString(s, ""))
}

// MaxDuration calculates how many seconds of the video will fit under TgUploadLimit.
// args: filesize (float64 in MB), duration (float64 in seconds)
// returns: int (seconds), error if calculation invalid
func MaxDuration(filesize, duration float64) (int, error) {
	if filesize <= 0 || duration <= 0 {
		return 0, fmt.Errorf("%w: filesize=%.2f MB, duration=%.2f s", myerrors.ErrInvalidInput, filesize, duration)
	}

	// Calculate ratio
	ratio := TgUploadLimit * FileSizeFix / filesize
	if ratio >= 1.0 {
		// Whole video fits
		return int(math.Floor(duration)), nil
	}

	// Calculate maximum allowed seconds
	maxSeconds := duration * ratio
	if maxSeconds < 1 {
		return 0, fmt.Errorf("%w too short: %.2f s", myerrors.ErrCalculatedDuration, maxSeconds)
	}

	return int(math.Floor(maxSeconds)), nil
}

func ConvertSecondsToMinSec(seconds int) string {
	minutes := seconds / SecondsInMinute
	seconds %= SecondsInMinute

	return strconv.Itoa(minutes) + ":" + strconv.Itoa(seconds)
}
