package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/slack-go/slack"
)

//SlashHandler - handles all slash commands registered for this app
func SlashHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("[INFO] Got slash command")
	jsonStr := ValidateRequest(w, r)
	if jsonStr == "ERROR" {
		log.Println("[ERROR] could not unmarshall the JSON request body")
	}
	command := slack.SlashCommand{}
	if err := json.Unmarshal([]byte(jsonStr), &command); err != nil {
		log.Printf("[ERROR] Failed to decode json message from slack: %s", jsonStr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var out bytes.Buffer
	var stderr bytes.Buffer

	switch command.Command {
	case "/createuser":
		log.Printf("[INFO] Got [%s] with arg [%s]", command.Command, command.Text)
		w.WriteHeader(http.StatusOK)
		cmd := exec.Command("aws", "iam", "create-user", "--user-name", command.Text)
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			log.Printf("[ERROR] Could not create AWS user\n", err)
		}
		type UserCreated struct {
			User struct {
				UserName string `json:"UserName"`
				Path string `json:"Path"`
				CreateDate string `json:"CreateDate"`
				UserID string `json:"UserId"`
				Arn string `json:"Arn"`
			}
		}
		var userCreated UserCreated
		if err := json.Unmarshal([]byte(out.String()), &userCreated); err != nil {
			log.Printf("[ERROR] Failed unmarshalling the create-user response")
			return
		}
		log.Printf("[DEBUG] User [%s] created successfully, ID [%s]", userCreated.User.UserName, userCreated.User.UserID )
		cmd = exec.Command("aws","iam","create-login-profile","--user-name", userCreated.User.UserName, "--password", "Ch@ng3Me!", "--password-reset-required")
		err = cmd.Run()
		if err != nil {
			log.Printf("[ERROR] Failed setting initial password for user [%s]", userCreated.User.UserName)
			return
		}
		log.Printf("[INFO] Succesfully set initial password for User [%s]", userCreated.User.UserName)
		str := fmt.Sprintf(`{"text":"User [%s] created successfully, and granted AWS Console access"}`, userCreated.User.UserName)
		jsonResponse := []byte(str)
		req, err := http.NewRequest("POST", command.ResponseURL, bytes.NewBuffer(jsonResponse))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("[ERROR] could not send verification to user",err)
		}
		defer resp.Body.Close()
		return
	case "/getrunningec2":
		log.Printf("[INFO] Got [%s] with arg [%s]", command.Command, command.Text)
		regions := GetEc2Regions()
		cmd := exec.Command("aws","ec2","describe-instances","--filters", "Name=instance-state-name,Values=running", "--query", "Reservations[*].Instances[*].{Instance:InstanceId,Key:KeyName}","--output","text", "--region","eu-west-2")
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			log.Printf("[ERROR] Failed getting running instances\n%s\n%s",err,stderr.String())
			return
		}
		fixed := strings.ReplaceAll(out.String(), "\n", "\\n")
		str := fmt.Sprintf("{\"text\":\"```%s\n\nREGIONS\n%s```\"}", fixed, regions)
		log.Printf("[DEBUG] sending response to slack\n%s", str)
		
		jsonResponse := []byte(str)
		req, err := http.NewRequest("POST", command.ResponseURL, bytes.NewBuffer(jsonResponse))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("[ERROR] could not send running instances to user",err)
		}
		defer resp.Body.Close()
		return
	default:
		log.Printf("[DEBUG] UNKNON COMMAND [%S]", command.Command)
	} 
}