package downloader

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/wader/goutubedl"

	c "github.com/ar2rworld/golang-telegram-video-downloader/app/constants"
	"github.com/ar2rworld/golang-telegram-video-downloader/app/logger"
	"github.com/ar2rworld/golang-telegram-video-downloader/app/myerrors"
	"github.com/ar2rworld/golang-telegram-video-downloader/app/platform"
	"github.com/ar2rworld/golang-telegram-video-downloader/app/utils"
)

type Downloader struct {
	Logger    *logger.Logger
	YtdlpPath string
}

func NewDownloader(l *logger.Logger, y string) *Downloader {
	goutubedl.Path = y
	return &Downloader{
		Logger: l,
		YtdlpPath: y,
	}
}

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
func (d *Downloader) DownloadVideo(ctx context.Context, url string, opts goutubedl.Options, do *goutubedl.DownloadOptions, prms *Parameters) (string, error) { //nolint: gocyclo,cyclop,funlen
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
	if err != nil && errors.Is(err, myerrors.ErrNoSizeInfo) && isDefaultSection {
		opts.DownloadSections = c.DefaultSections
	}

	if needsCutting {
		sections, err := prms.Platform.MaxDuration(&result)
		if err != nil {
			d.Logger.Error().Err(err).Msg("error max duration: ")
			sections = c.DefaultSections
		}
		opts.DownloadSections = sections
	}

	result, err = goutubedl.New(ctx, url, opts)
	if err != nil {
		return "", ReturnNewRequestError(err)
	}

	// TODO: fix downloading audio with goutubedl
	if do.DownloadAudioOnly {
		filename := result.Info.Title + "." + result.Info.Ext
		filename = cutString(filename, c.MaxFileNameLength)
		filename = ConvertToUTF8(filename)
		filename = RemoveNonAlphanumericRegex(filename)
		filename = strings.ReplaceAll(filename, " ", "")
		filename = path.Join(os.TempDir(), filename)
		prms.AddTempFile(filename)
		return d.DownloadAudio(ctx, url, opts.Cookies, filename)
	}

	d.Logger.Info().Str("platform", prms.Platform.Name()).
		Str("section", opts.DownloadSections).
		Float64("filesize", result.Info.Filesize).
		Float64("filesize_approx", result.Info.FilesizeApprox).
		Msg("DownloadWithOptions")

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
	d.Logger.Info().Float64("filesize", s).Msg("Filesize of downloaded video")

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

func (d *Downloader) DownloadAudio(ctx context.Context, url, cookies, filename string) (string, error) {
	args := []string{
		"-o", filename,
		"--extract-audio",
		"--print", "after_move:filepath",
	}

	if cookies != "" {
		args = append(args, "--cookies", cookies)
	}

	args = append(args, url)

	cmd := exec.CommandContext(ctx, d.YtdlpPath, args...)

	output, err := cmd.CombinedOutput()
	outStr := string(output)

	if err != nil {
		return "", fmt.Errorf("yt-dlp error: %v\n%s", err, outStr)
	}

	// Extract the last non-empty line (usually the file path)
	lines := strings.Split(outStr, "\n")
	var finalPath string
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if line != "" {
			finalPath = line
			break
		}
	}

	if finalPath == "" {
		finalPath = filename
	}

	d.Logger.Info().Str("finalPath", finalPath).Msg("download audio")
	return finalPath, nil
}

