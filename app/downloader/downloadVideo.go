package downloader

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/wader/goutubedl"

	c "github.com/ar2rworld/golang-telegram-video-downloader/app/constants"
	"github.com/ar2rworld/golang-telegram-video-downloader/app/myerrors"
	"github.com/ar2rworld/golang-telegram-video-downloader/app/platform"
	"github.com/ar2rworld/golang-telegram-video-downloader/app/utils"
)

var YtdlpPath = "yt-dlp_macos" //nolint:gochecknoglobals

const FileSizeFix = 0.5

type Parameters struct {
	Platform platform.Platform
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
	goutubedl.Path = YtdlpPath

	isDefaultSection := opts.DownloadSections == c.DefaultSections
	// if DefaultSections is set, select video section under TgUploadLimit
	if isDefaultSection {
		opts.DownloadSections = ""
	}

	result, err := goutubedl.New(ctx, url, opts)
	if err != nil {
		return "", ReturnNewRequestError(err)
	}

	if do == nil {
		do = &goutubedl.DownloadOptions{}
	}

	filter, err := prms.Platform.SelectFormat(result.Formats())
	if err != nil && !errors.Is(err, myerrors.ErrNoSuitableFormat) {
		return "", err
	}
	do.Filter = filter

	needsCutting, err := prms.Platform.NeedCut(&result)
	if err != nil && errors.Is(err, myerrors.ErrNoSizeInfo) {
		opts.DownloadSections = c.DefaultSections
	}

	if needsCutting {
		sections, err := prms.Platform.MaxDuration(&result)
		if err != nil {
			log.Println("*** error max duration: ", err)
			sections = c.DefaultSections
		}
		opts.DownloadSections = sections
	}

	result, err = goutubedl.New(ctx, url, opts)
	if err != nil {
		return "", ReturnNewRequestError(err)
	}

	// Weird stuff to make yt-dlp download only audio
	if !do.DownloadAudioOnly && do.Filter == "" {
		do.Filter = "best"
	}

	log.Printf("*** DownloadWithOptions: section: %s, filesize: %f, filesize approx: %f\n", opts.DownloadSections, result.Info.Filesize, result.Info.FilesizeApprox)

	downloadResult, err := result.DownloadWithOptions(ctx, *do)
	if err != nil {
		return "", err
	}
	defer downloadResult.Close()

	title := RemoveNonAlphanumericRegex(ConvertToUTF8(result.Info.Title))
	output := cutString(title, c.MaxFileNameLength) + result.Info.Ext

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
	if prms.Platform.RemuxRequired() {
		filename, err = RemuxToMP4(ctx, filename)
		if err != nil {
			return "", fmt.Errorf("could not remux video to mp4: %w", err)
		}

		prms.AddTempFile(filename)
	}

	s, _ := FileSizeMB(filename)
	log.Printf("*** Filesize of downloaded video: %f\n", s)

	return filename, nil
}

// FileSizeMB takes a filepath and returns file size in MB (float64)
func FileSizeMB(filepath string) (float64, error) {
	fileInfo, err := os.Stat(filepath)
	if err != nil {
		return 0, err
	}

	sizeBytes := fileInfo.Size()                                 // size in bytes
	sizeMB := utils.BytesToMb(float64(sizeBytes)) // convert to MB

	return sizeMB, nil
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


func ReturnNewRequestError(err error) error {
	if strings.Contains(err.Error(), myerrors.UnsupportedURL) {
		return myerrors.ErrUnsupportedURL
	}

	if strings.Contains(err.Error(), myerrors.VideoUnavailable) {
		return myerrors.ErrVideoUnavailable
	}

	if strings.Contains(err.Error(), myerrors.RequestedContentIsNotAvailable) {
		return myerrors.ErrRequestedContentIsNotAvailable
	}

	return err
}
