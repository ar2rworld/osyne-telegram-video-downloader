package constants

type Options struct {
	Prod                 bool   `description:"Prod env"                env:"PROD"                   long:"prod"                   short:"p"`
	YtdlpPath            string `description:"YT-DLP path"             env:"YT_DLP_PATH"            long:"ytdlp_path"             required:"true" short:"d"`
	BotToken             string `description:"Telegram bot token"      env:"BOT_TOKEN"              long:"bot_token"              required:"true" short:"t"`
	AdminID              int64  `description:"Telegram admin id"       env:"ADMIN_ID"               long:"admin_id"               required:"true" short:"a"`
	LogChannelID         int64  `description:"Telegram log channel id" env:"LOG_CHANNEL_ID"         long:"log_channel_id"         required:"true" short:"l"`
	CookiesPath          string `description:"Cookies path"            env:"COOKIES_PATH"           long:"cookies_path"           required:"true" short:"c"`
	InstagramCookiesPath string `description:"Instagram cookies path"  env:"INSTAGRAM_COOKIES_PATH" long:"instagram_cookies_path" required:"true" short:"i"`
	YouTubeCookiesPath   string `description:"YouTube cookies path"    env:"GOOGLE_COOKIES_PATH"    long:"youtube_cookies_path"   required:"true" short:"y"`
}
