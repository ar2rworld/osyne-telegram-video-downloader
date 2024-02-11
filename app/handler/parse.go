package handler

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// If user specifies -s with some argument to download only section of video
func parse(s string) (string, error) {
	fs := flag.NewFlagSet("sections", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	var sections string

	fs.StringVar(&sections, "s", fmt.Sprintf("*0:0-0:%d", Duration), "Download sections")

	err := fs.Parse(strings.Split(s, " "))
	if err != nil {
		fmt.Println("Error parsing flags:", err)
		return "", err
	}

	return sections, nil
}
