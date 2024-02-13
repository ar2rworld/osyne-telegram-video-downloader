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
		fmt.Println("Error parsing flags:", err)
		return DefaultSections, err
	}

	return sections, nil
}
