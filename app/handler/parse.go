package handler

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/jessevdk/go-flags"
)

const DefaultSections = "*0:0-0:30"

type UserOptions struct {
	Sections     *string `description:"Download sections"        long:"sections"      short:"s"`
	ExtractAudio *bool   `description:"Extract audio from video" long:"extract-audio" short:"x"`
}

// parses the string and return UserOptions
func parse(s string) (*UserOptions, error) {
	opts := &UserOptions{}

	parser := flags.NewParser(opts, flags.None)

	args := strings.Split(s, " ")
	_, err := parser.ParseArgs(args)
	if err != nil {
		return nil, fmt.Errorf("error parsing arguments: %w", err)
	}

	return opts, nil
}

// parse youtube url for current time argument
func parseCurrentTime(videoURL string) string {
	// Claude 3.5 Sonnet
	parsedURL, err := url.Parse(videoURL)
	if err != nil {
		return ""
	}

	// Check query parameters
	if t := parsedURL.Query().Get("t"); t != "" {
		// Remove any non-digit characters and parse as int
		re := regexp.MustCompile(`\d+`)
		if match := re.FindString(t); match != "" {
			if _, err := strconv.Atoi(match); err == nil {
				return match
			}
		}
	}

	return ""
}
