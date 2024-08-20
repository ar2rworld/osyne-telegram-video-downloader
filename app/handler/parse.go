package handler

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const DefaultSections = "*0:0-0:30"

// If user specifies -s with some argument to download only section of video
func parse(s string) (string, error) {
	fs := flag.NewFlagSet("sections", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	var sections string

	fs.StringVar(&sections, "s", DefaultSections, "Download sections")

	err := fs.Parse(strings.Split(s, " "))
	if err != nil {
		return "", fmt.Errorf("error parsing sections: %w", err)
	}
	// if flag parsed without errors but sections wasnot provided
	if sections == "" {
		sections = DefaultSections
	}

	return sections, nil
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
