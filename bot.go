package main

import (
	"log"
	"net/http"
	"os"

	"github.com/slack-go/slack"
)

// Load environment vars (Path to SSL certificate and key, the bot token, and the virfication token [to verify incoming messages come from Slack - this is depricated will switch to `signing secret` sometime])
//TODO: Switch to signing secret - rather than the, depricated, verification token
var certPath, keyPath, botToken, verificationToken, awsKey, awsSecret = getEnvVars()

// Start the Slack client
var api = slack.New(botToken)

func main() {
	checkVars()
	// listen for Slack events here
	http.HandleFunc("/event", EventHandler)
	// listen for Slack user interactions here
	http.HandleFunc("/interact", InteractionHandler)
	// listen for slash command `createUser`
	http.HandleFunc("/slash", SlashHandler)
	//Start listening
	log.Println("[INFO] Starting to listen on port 443")
	err := http.ListenAndServeTLS(":443", certPath, keyPath, nil)
	if err != nil {
		log.Fatal("[ERROR] ListenAndServeTLS:", err)
	}
}

func checkVars() {
	if certPath == "" || keyPath == "" {
		log.Println("[ERROR] SSL configuration error. Missing environment variables (either SSL_CERT_PATH or SSL_KEY_PATH)")
	}
	if botToken == "" {
		log.Println("[ERROR] Missing bot token. Missing environment variable (SLACK_BOT_TOKEN)")
	}
	if verificationToken == "" {
		log.Println("[ERROR] Missing verification token. Missing environment variable (SLACK_VERIFICATION_TOKEN)")
	}
	if awsKey == "" {
		log.Println("[ERROR] missing AWS_KEY_ID environment variable")
	}
	if awsSecret == "" {
		log.Println("[ERROR] missing AWS_SECRET_KEY environment variable")
	}
}

func getEnvVars() (certPath string, keyPath string, botToken string, verificationToken string, awsKey string, awsSecret string) {
	certPath = os.Getenv("SSL_CERT_PATH")
	keyPath = os.Getenv("SSL_KEY_PATH")
	botToken = os.Getenv("SLACK_BOT_TOKEN")
	verificationToken = os.Getenv("SLACK_VERIFICATION_TOKEN")
	awsKey = os.Getenv("AWS_KEY_ID")
	awsSecret = os.Getenv("AWS_SECRET_KEY")
	return
}
