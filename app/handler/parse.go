package handler

import (
	"flag"
	"fmt"
	"os"
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
