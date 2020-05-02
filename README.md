# This is a Slack bot written in Golang built to communicate with AWS
## The bot is packaged in a GO module

### Prerequisites
1. Create a server that will run your bot (This server needs to run in your VPC and have *******permissions*******)
   * The server needs to be accessible by Slack
   * The server needs to have a valid SSL ceritificate (I used LetsEncrypt for testing)
2. Create a new Slack app (https://api.slack.com/apps)
   * **During the authentication phase of the app creation you will be asked to reply to a specific request with a specific response to validate ownership of the webserver.**
3. Subscribe to the following events (Using the Events API functionality) and set the Request URL to https://<yourServer>/event
   * `app_mention` (So the bot can answer when mentioned)
   * `memeber_joined_channel` (So the bot can greet new members)
4. Add Interactivity functionality and set the Request URL to https://\<youServer\>/interact
5. Register the following 2 Slash Commands to be send to https://\<youServer\>/slash
   1. `/createUser \<userName\>`
   2. `/getrunningec2`

### To run this bot:
1. Install Golang
2. Create, build and install the bot module
    ```cd bot
    go mod init <moduleName (i.e acme/slackbot)>
    go install <modulName ^>
3. Set the following environement vars:
* SSL_CERT_PATH (SSL cert / bundle)
* SSL_KEY_PATH (SSL private key)
* SLACK_BOT_TOKEN (Slack bot token (starting with "xoxb-"))
* SLACK_VERIFICATION_TOKEN (Slack verification token (this is the **old** token - to verify request source))
* AWS_KEY_ID
* AWS_SECRET_KEY
4. Run the bot
   `sudo -E go run <moduleName from step 2> &`