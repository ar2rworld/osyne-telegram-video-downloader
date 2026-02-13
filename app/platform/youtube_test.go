package platform

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wader/goutubedl"

	c "github.com/ar2rworld/golang-telegram-video-downloader/app/constants"
	"github.com/ar2rworld/golang-telegram-video-downloader/app/myerrors"
)

func TestSelectFormat(t *testing.T) { //nolint:funlen
	yt := NewYoutube("")

	tests := []struct {
		name       string
		formats    []goutubedl.Format
		wantFilter string
		wantErr    error
	}{
		{
			name:    "no formats",
			wantErr: myerrors.ErrNoSuitableFormat,
		},
		{
			name: "formats below min resolution",
			formats: []goutubedl.Format{
				{FormatID: "1", Width: 100, Height: 100, VCodec: "avc1.42", ACodec: "none", Filesize: 1024 * 1024},
			},
			wantErr: myerrors.ErrNoSuitableFormat,
		},
		{
			// NOTE: selectBestFormat compares raw Filesize against TgUploadLimit (50)
			// without converting bytes to MB — a known inconsistency with the
			// brute-force path which correctly uses BytesToMb.
			name: "complete format under limit",
			formats: []goutubedl.Format{
				{
					FormatID: "22", Width: 1280, Height: 720,
					VCodec: c.VideoCodec + ".42", ACodec: c.AudioCodec + ".40",
					Filesize: 30, FilesizeApprox: 30,
				},
			},
			wantFilter: "22",
		},
		{
			name: "video+audio combo under limit",
			formats: []goutubedl.Format{
				{
					FormatID: "137", Width: 1920, Height: 1080,
					VCodec: c.VideoCodec + ".42", ACodec: "none",
					Filesize: 30 * c.BytesInKByte * c.BytesInKByte,
				},
				{
					FormatID: "140",
					VCodec:   "none", ACodec: c.AudioCodec + ".40",
					Filesize: 5 * c.BytesInKByte * c.BytesInKByte,
				},
			},
			wantFilter: "137+140",
		},
		{
			name: "all formats exceed limit",
			formats: []goutubedl.Format{
				{
					FormatID: "137", Width: 1920, Height: 1080,
					VCodec: c.VideoCodec + ".42", ACodec: "none",
					Filesize: 40 * c.BytesInKByte * c.BytesInKByte,
				},
				{
					FormatID: "140",
					VCodec:   "none", ACodec: c.AudioCodec + ".40",
					Filesize: 20 * c.BytesInKByte * c.BytesInKByte,
				},
			},
			wantErr: myerrors.ErrNoSuitableFormat,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter, err := yt.SelectFormat(tt.formats)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				assert.Empty(t, filter)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantFilter, filter)
		})
	}
}
