package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)
//EventHandler - recieves the Slack events that the app is subscribed for
func EventHandler(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	body := buf.String()
	// If you don't need source authentication uncomment the next line, and comment the one after it
	// eventsAPIEvent, e := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
	// and disable the environment variable check in `checkVars()`
	eventsAPIEvent, e := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: verificationToken}))
	// eventsAPIEvent, e := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionVerifyToken()
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
	}
	log.Println("Received event from Slack")
	if eventsAPIEvent.Type == slackevents.URLVerification {
		log.Println("------ URL Verification Event")
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(body), &r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "text")
		w.Write([]byte(r.Challenge))
	}
	if eventsAPIEvent.Type == slackevents.CallbackEvent {
		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			log.Printf("------------AppMentionEvent from user ID [%s]", ev.User)
			msg := fmt.Sprintf("Hi <@%s>, Here's what I can do:\n_/createUser <username>_ Will create an AWS user\n_/getrunningec2_ will show a list of running servers (currently in Lodon - I'll soon be upgraded to see a bit more of the world :eyeglasses:",ev.User)
			_, _, err := api.PostMessage(ev.Channel, slack.MsgOptionText(msg, false))
			if err != nil {
				log.Printf("[ERROR] Error sending message to [%s] in channel [%s]\n%s", ev.User, ev.Channel, err)
			}
		case *slackevents.MemberJoinedChannelEvent:
			log.Println("-------------MemberJoinedChannelEvent")
			msg := fmt.Sprintf("Hi, <@%s>, Welcome! You may mention me in a message to see my capabilities", ev.User)
			_, _, err := api.PostMessage(ev.Channel, slack.MsgOptionText(msg, false))
			if err != nil {
				log.Printf("[ERROR] Error sending message to [%s] in channel [%s]\n%s", ev.User, ev.Channel, err)
			}
		case *slackevents.MessageActionResponse:
			log.Println("-------------MessageActionResponse")

		}
	}
}