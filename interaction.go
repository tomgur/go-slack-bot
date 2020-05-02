package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/slack-go/slack"
)

//InteractionHandler - recieves all the app's interactions
func InteractionHandler(w http.ResponseWriter, r *http.Request) {
	jsonStr := ValidateRequest(w, r)
	if jsonStr == "ERROR" {
		log.Printf("[ERROR] could not unmarshall the JSON request body")
	}
	var message slack.InteractionCallback
	if err := json.Unmarshal([]byte(jsonStr), &message); err != nil {
		log.Printf("[ERROR] Failed to decode json message from slack: %s", jsonStr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("[DEBUG] - unmarshalled the request body.\n%+v", message)
	switch message.Type {
	case slack.InteractionTypeInteractionMessage:
		log.Printf("Got response from user [%s]", message.User.Name)
		log.Printf("[DEBUG] CallbackID %+v", message.Message)
		log.Printf("[DEBUG] `message.OriginalMessage`\n%+v", message.OriginalMessage)
		actions := message.ActionCallback.AttachmentActions
		action := actions[0]
		switch action.Name {
		case "actionSelect":
			log.Println("[DEBUG] User selected an action")
			value := action.SelectedOptions[0].Value
			// Overwrite original drop down message.
			originalMessage := message.OriginalMessage
			originalMessage.Attachments[0].Text = fmt.Sprintf("OK to order %s ?", strings.Title(value))
			originalMessage.Attachments[0].Actions = []slack.AttachmentAction{
				{
					Name:  "actionStart",
					Text:  "Yes",
					Type:  "button",
					Value: "start",
					Style: "primary",
				},
				{
					Name:  "actionCancel",
					Text:  "No",
					Type:  "button",
					Style: "danger",
				},
			}
			w.Header().Add("Content-type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(&originalMessage)
			return
		case "actionStart":
			log.Println("[DEBUG] User approved selection")
			title := ":ok: your order was submitted! yay!"
			log.Printf("Going to respond to message: %+v", message.Message)
			responseMessage(w, message.OriginalMessage, title, "")
			return
		case "actionCancel":
			log.Println("[DEBUG] User canceled action")
			title := fmt.Sprintf("@%s canceled the request", message.User.Name)
			responseMessage(w, message.Message, title, "")
			return
		default:
			log.Printf("[ERROR] ]Invalid action was submitted: %s", action.Name)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func responseMessage(w http.ResponseWriter, original slack.Message, title, value string) {
	log.Println("[DEBUG] title: ", title)
	log.Println("[DEBUG] original: ", original)
	original.Msg.Attachments[0].Actions = []slack.AttachmentAction{} // empty buttons
	original.Msg.Attachments[0].Fields = []slack.AttachmentField{
		{
			Title: title,
			Value: value,
			Short: false,
		},
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&original)
}