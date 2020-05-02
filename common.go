package main

import (
	"log"
	"net/http"
	"strings"
	"net/url"
	"io/ioutil"
	"bytes"
	"os/exec"

)

// ValidateRequest - Validates that the request is POST, and transforms the body from URL encoding, to single-line JSON
func ValidateRequest(w http.ResponseWriter, r *http.Request) string {
	log.Println("[DEBUG] Validating request")
	if r.Method != http.MethodPost {
		log.Printf("[ERROR] Invalid method: %s", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return "ERROR"
	}
	
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("[ERROR] Failed to read request body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return "ERROR"
	}
	jsonStr, err := url.QueryUnescape(string(buf))
	if err != nil {
		log.Printf("[ERROR] Failed to unescape request body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return "ERROR"
	}
	fixed := "{\"" + jsonStr + "\"}"
	fixed = strings.ReplaceAll(fixed, "&", "\",\"")
	result := strings.ReplaceAll(fixed, "=", "\":\"")
	log.Println("[DEBUG] Request Validated")
	return result
}

//GetEc2Regions - returns a string of regions from the AWS CLI
func GetEc2Regions() string {
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("aws", "ec2", "describe-regions", "--output", "text")
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Printf("[ERROR] Could not get AWS regions\n%s\n%s", out, err)
	}
	return out.String()
}