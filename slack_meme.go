package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

// https://api.slack.com/incoming-webhooks
type slackRequest struct {
	Text        string       `json:"text"`
	Username    string       `json:"username"`
	Channel     string       `json:"channel"`
	Attachments []Attachment `json:"attachments"`
}

// https://api.slack.com/docs/attachments
type Attachment struct {
	Fallback   string `json:"fallback,omitempty"`
	Color      string `json:"color, omitempty"`
	PreText    string `json:"pretext,omitempty"`
	AuthorName string `json:"author_name,omitempty"`
	AuthorLink string `json:"author_link,omitempty"`
	AuthorIcon string `json:"author_icon,omitempty"`
	Title      string `json:"title,omitempty"`
	TitleLink  string `json:"title_link,omitempty"`
	Text       string `json:"text,omitempty"`
	ImageURL   string `json:"image_url,omitempty"`
	ThumbURL   string `json:"thumb_url,omitempty"`
}

var (
	regEx = regexp.MustCompile(`(.*[\s])(top:)(.*[\s]+)(bottom:)(.*)`)
)

// Read the incoming request, then send the response
func readRequest(r *http.Request) string {
	err := r.ParseForm()
	// TODO: Change HTTP status code
	if err != nil {
		return string(err.Error())
	}
	// Check incoming token against config
	if len(r.Form["token"]) == 0 || r.Form["token"][0] != os.Getenv("SLACK_TOKEN") {
		return "Incoming token does not match configured SLACK_TOKEN"
	}
	if len(r.Form["text"]) == 0 {
		return "Empty Text"
	}
	text := strings.TrimSpace(r.Form["text"][0])
	if text == "list" {
		return "https://api.imgflip.com/popular_meme_ids"
	}
	parsedText := regEx.FindStringSubmatch(text)

	if parsedText == nil {
		return "Try again: /meme 61579 top: ONE DOES NOT SIMPLY bottom: CREATE A MEME"
	}
	templateId := parsedText[0]
	topText := parsedText[3]
	bottomText := parsedText[5]

	channel := r.Form["channel_id"][0]
	memeUrl, err := createMeme(templateId, topText, bottomText)
	if err != nil {
		return "Failed to create meme"
	}
	err = sendMeme(channel, memeUrl)
	if err != nil {
		return "Failed to send message."
	}
	return fmt.Sprintf("Meme created : [%s]", memeUrl)
}

// Compose imgflip meme
func createMeme(templateId, topText, bottomText string) (string, error) {
	values := url.Values{}
	values.Set("template_id", templateId)
	values.Set("username", os.Getenv("IMGFLIP_USERNAME"))
	values.Set("password", os.Getenv("IMGFLIP_PASSWORD"))
	values.Set("text0", topText)
	values.Set("text1", bottomText)
	resp, err := http.PostForm("https://api.imgflip.com/caption_image", values)

	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", err
	}

	if !data["success"].(bool) {
		return "", errors.New(data["error_message"].(string))
	}

	url := data["data"].(map[string]interface{})["url"].(string)

	return url, nil
}

// Send meme back to the channel
func sendMeme(channel, memeUrl string) error {
	url := os.Getenv("SLACK_WEBHOOK")
	payload, err := json.Marshal(slackRequest{
		Channel:     channel,
		Username:    "SlackMeme",
		Attachments: []Attachment{Attachment{ImageURL: memeUrl}},
	})
	if err != nil {
		return err
	}
	_, err = http.Post(url, "application/json", bytes.NewBuffer(payload))
	return err
}

func main() {
	rand.Seed(time.Now().UnixNano())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		result := readRequest(r)
		fmt.Fprintf(w, result)
	})
	http.ListenAndServe(fmt.Sprintf(":%d", 5000), nil)
}
