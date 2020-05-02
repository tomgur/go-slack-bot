This is a Slack bot written in Golang.
======================================

To run this bot:
----------------
1. Install Golang
2. Build the bot
> cd bot
> go build bot.go
3. Set the following environement vars:
* SSL_CERT_PATH (SSL cert / bundle)
* SSL_KEY_PATH (SSL private key)
* SLACK_BOT_TOKEN (Slack bot token (starting with "xoxb-"))
* SLACK_VERIFICATION_TOKEN (Slack verification token (this is the **old** token - to verify request source))
4. Run the bot
> ./bot (bot.exe on Windows)