module github.com/ar2rworld/golang-telegram-video-downloader

go 1.24.0

require (
	github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1
	github.com/stretchr/testify v1.9.0
	github.com/wader/goutubedl v0.0.0-20260211162955-2c534af3ada4
)

require golang.org/x/sys v0.41.0 // indirect

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/jessevdk/go-flags v1.6.1
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/go-telegram-bot-api/telegram-bot-api/v5 => github.com/ar2rworld/telegram-bot-api/v5 v5.0.0
