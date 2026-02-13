package downloader

import (
	"errors"
	"testing"

	"github.com/wader/goutubedl"
)

func TestSelectFormat_NoSuitableFormat(t *testing.T) {
	// Empty formats should return ErrNoSuitableFormat
	filter, ext, err := SelectFormat([]goutubedl.Format{})
	if !errors.Is(err, ErrNoSuitableFormat) {
		t.Fatalf("expected ErrNoSuitableFormat, got %v", err)
	}

	if filter != "" {
		t.Errorf("expected empty filter, got %q", filter)
	}

	if ext != "" {
		t.Errorf("expected empty ext, got %q", ext)
	}
}

func TestSelectFormat_FilterNotOverwrittenOnError(t *testing.T) {
	// Simulates the fixed logic in DownloadVideo: when SelectFormat returns
	// ErrNoSuitableFormat, do.Filter must remain FilterBest and not be overwritten.
	do := &goutubedl.DownloadOptions{}

	filter, e, err := SelectFormat([]goutubedl.Format{})
	if err != nil && !errors.Is(err, ErrNoSuitableFormat) {
		t.Fatalf("unexpected error: %v", err)
	}

	if err != nil {
		do.Filter = FilterBest
	} else {
		do.Filter = filter
		_ = e
	}

	if do.Filter != FilterBest {
		t.Errorf("expected do.Filter to be %q, got %q", FilterBest, do.Filter)
	}
}

func TestSelectFormat_ValidFormats(t *testing.T) {
	formats := []goutubedl.Format{
		{
			FormatID: "v1",
			Ext:      ExtMP4,
			Width:    1920,
			Height:   1080,
			VCodec:   "avc1.something",
			ACodec:   "none",
			Filesize: 10 * BytesInKByte * BytesInKByte, // 10 MB
		},
		{
			FormatID: "a1",
			Width:    0,
			Height:   0,
			VCodec:   "none",
			ACodec:   "mp4a.something",
			Filesize: 5 * BytesInKByte * BytesInKByte, // 5 MB
		},
	}

	filter, ext, err := SelectFormat(formats)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if filter != "v1+a1" {
		t.Errorf("expected filter %q, got %q", "v1+a1", filter)
	}

	if ext != ExtMP4 {
		t.Errorf("expected ext %q, got %q", ExtMP4, ext)
	}
}
