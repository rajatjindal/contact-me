package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	spinhttp "github.com/fermyon/spin-go-sdk/http"
	"github.com/fermyon/spin-go-sdk/variables"
)

func init() {
	spinhttp.Handle(func(w http.ResponseWriter, r *http.Request) {
		cors(w)

		if r.Method == http.MethodOptions {
			return
		}

		raw, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var contact Contact
		err = json.Unmarshal(raw, &contact)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = sendMsg(r.Context(), contact)
		if err != nil {
			fmt.Println("ERROR: failed to send contact request", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	})
}

func cors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "authorization, content-type")
}

type Contact struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Msg   string `json:"msg"`
}

func sendMsg(ctx context.Context, contact Contact) error {
	payload := fmt.Sprintf(`{
	"blocks": [
		{
			"type": "section",
			"text": {
				"type": "mrkdwn",
				"text": "New contact request from rajatjindal.com:"
			}
		},
		{
			"type": "section",
			"text": {
				"type": "mrkdwn",
				"text": "*Name*: %s"
			}
		},
		{
			"type": "section",
			"text": {
				"type": "mrkdwn",
				"text": "*Msg*: %s"
			}
		},
		{
			"type": "section",
			"text": {
				"type": "mrkdwn",
				"text": "*Email*: %s"
			}
		},
	]
}`, contact.Name, contact.Msg, contact.Email)

	slackurl, err := variables.Get("slack_webhook")
	if err != nil {
		return err
	}

	resp, err := spinhttp.Post(slackurl, "application/json", strings.NewReader(payload))
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send contact request. Expected Status code: %d, actual: %d", http.StatusOK, resp.StatusCode)
	}

	return nil
}

func main() {}
