# OsyNe telegram video downloader

**OsyNe telegram video downloader** is a Go-based Telegram bot that allows users to download videos from various sources such as TikTok, Twitter, and Instagram. It downloads the video content and sends it back to the user as a reply to their message.

## Features

- Download videos from TikTok, Twitter, and Instagram
- Sends downloaded videos back to users as a reply
- Provides support for downloading videos from Instagram with authentication (using cookies(depricated))

## Getting Started

### Prerequisites

1. Go programming language installed on your system.
2. A Telegram Bot API token. You can obtain it by creating a new bot on Telegram via the BotFather.
3. A file containing Instagram **Netscape formatted** cookies (required if you want to download videos from Instagram). Learn how to obtain these cookies in the section below.

### Installation

1. Clone this repository to your local machine:

```
git clone https://github.com/ar2rworld/telegram-video-downloader-bot.git
cd telegram-video-downloader-bot
```

2. Install the required dependencies:

```
go get github.com/go-telegram-bot-api/telegram-bot-api/v5
```

3. Complete installation process for [yt-dlp](https://github.com/yt-dlp/yt-dlp)

### Usage

1. Set the required environment variables:

    - `BOT_TOKEN`: Your Telegram Bot API token.
    - `ADMIN_ID`: ID of the bot admin (the bot will send a "Hello, boss" message to this ID when it starts).

2. Build and run the bot:

```
go build
./telegram-video-downloader-bot
```

3. Interact with the bot:

    - Send a video URL from TikTok, Twitter, or Instagram to the bot.
    - The bot will download the video and send it back as a reply.

## License

This project is licensed under the [MIT License](LICENSE).

## TODO:
- CICD
  - push to the docker hub
- Feature: Delete text of the requesting user and leave only video
- Youtube:
  - limit with 1 minute
  - Video part specified

*Thanks to ChatGPT3.5*