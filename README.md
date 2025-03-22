# OsyNe telegram video downloader

**OsyNe telegram video downloader** is a Go-based Telegram bot that allows users to download videos from various sources such as TikTok, Twitter, Youtube shorts and Instagram. It downloads the video content and sends it back to the user as a reply to their message.

## Features

- Download videos from TikTok, Twitter, Youtube shorts, and Instagram
- Sends downloaded videos back to users as a reply
- Provides support for downloading videos from Instagram with authentication (using cookies)

## Getting Started

### Prerequisites

1. Go programming language installed on your system.
2. A Telegram Bot API token. You can obtain it by creating a new bot on Telegram via the BotFather.
3. A file containing Instagram **Netscape formatted** cookies (required if you want to download videos from Instagram).

### Installation

1. Clone this repository to your local machine:

```
git clone https://github.com/ar2rworld/osyne-telegram-video-downloader.git
cd osyne-telegram-video-downloader
```

2. Install the required dependencies:

```
go mod download
```

3. Complete installation process for [yt-dlp](https://github.com/yt-dlp/yt-dlp)

4. Add environment variables for ADMIN_ID, BOT_TOKEN, and INSTAGRAM_COOKIES_FILE

5. Run application
```
go run ./app
```

### OR

Use docker-compose from [docker hub](https://hub.docker.com/r/ar2rworld/osyne-telegram-video-downloader)

1. Clone repo
```
git clone https://github.com/ar2rworld/osyne-telegram-video-downloader.git
cd osyne-telegram-video-downloader
```

2. Rename `template.env` to  `.env` file and add tokens

3. Run container
```
docker-compose --env-file .env up -d
```

### Usage

1. Set the required environment variables:

    - `BOT_TOKEN`: Your Telegram Bot API token.
    - `ADMIN_ID`: ID of the bot admin (the bot will send a "Hello, boss" message to this ID when it starts).
    - `INSTAGRAM_COOKIES_FILE`: Instagram session cookies file(REMOVED)
    - `ARTIFACTS_PATH`: Instagram cookie file path on the server to pass it to the container through the volume(REMOVED)

2. Build and run the bot:

```
go build -o ./osynetelegramvideodownloader ./app
./osynetelegramvideodownloader
```

3. Interact with the bot:

    - Send a video URL from TikTok, Twitter, Youtube shorts or Instagram to the bot.
    - The bot will download the video and send it back as a reply.
    - Send `osyndaisyn ba?` to chat and bot responds with `osyndaymyn` (In kazakh "Are you here?" and "I am here" - responce)

### Run linters
1. Install [pre-commit](https://pre-commit.com/#install)

2. Run `pre-commit run -a`
## License

This project is licensed under the [MIT License](LICENSE).

## TODO:
- CICD
  - push to the docker hub
    - (Done)
- Feature: Delete text of the requesting user and leave only video
  - (Bot's cannon react to messages or edit messages)
- Youtube:
  - limit with 1 minute
  - Video part specified
    - (Implemented with shorts)
- Instagram will require cookie settings
  - adding volume to the local folder with cookies.txt file
  - read a file everytime instagram.com requested
  - setting up cookies from env var (Done)
  - bot will replace cookies file in direct message for google and instagram
- Add https://golangci-lint.run/, locally could be run with https://pre-commit.com/, or with editor extention
  - https://gist.github.com/pantafive/3296201ef3dc14a71139cae157aa8c34
  - Also add to cicd (Done)
- Improve the functionality with Steam (downloaded steam -> telegram new video steam)
- Grafana
  - vector dev log collector
- Update the cleanUrl to add regexp finding the link
  - in the range : if the message has two links and the first one is not to tiktok
- Delete message after posting a video
  - add who sent, link to video


*Thanks to @pantafive, ChatGPT 3.5, 4o*
