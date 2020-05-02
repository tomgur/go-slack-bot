package main

import (
	"github.com/slack-go/slack"
)

// MainAttachment is the attachment sent when the bot is mentioned
var MainAttachment = slack.Attachment{
	Text:       "Which region? :world_map:",
	Color:      "#f9a41b",
	CallbackID: "region",
	Actions: []slack.AttachmentAction{
		{
			Name: "actionSelect",
			Type: "select",
			Options: []slack.AttachmentActionOption{
				{
					Text:  "North Verginia",
					Value: "us-east-1",
				},
				{
					Text:  "London",
					Value: "eu-west-2",
				},
			},
		},
		{
			Name:  "actionCancel",
			Text:  "Cancel",
			Type:  "button",
			Style: "danger",
		},
	},
}
