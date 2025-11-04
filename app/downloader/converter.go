package downloader

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
)

func Convert(ctx context.Context, path string) (string, error) {
	outputFile := path + ".mp4"

	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-y", "-i", path,
		"-c:v", "copy",
		"-c:a", "aac", "-b:a", "128k",
		outputFile,
	)

	cmd.Stderr = os.Stderr

	cmd.Stdout = os.Stdout

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return outputFile, nil
}

func RemuxToMP4(ctx context.Context, path string) (string, error) {
	// Change extension to .mp4
	outputFile := filepath.Base(path) + ".mp4"

	cmd := exec.CommandContext(ctx,
		"ffmpeg",
		"-y",       // overwrite without asking
		"-i", path, // input file
		"-c:v", "libx264", // copy video stream
		"-c:a", "copy", // copy audio stream
		"-movflags", "+faststart", // move mp4 metadata to beginning
		"-avoid_negative_ts", "make_zero", // adjust timestamps
		outputFile,
	)

	// Pipe ffmpeg output to current stderr/stdout
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return outputFile, nil
}
