package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

func handleWebhook(w http.ResponseWriter, r *http.Request) {

	hasher := hmac.New(sha1.New, []byte(os.Getenv("PUSHER_BEAMS_WEBHOOK_SECRET")))

	type response struct {
		Message string `json:"message"`
		Status  bool   `json:"status"`
	}

	if r.Header.Get("Webhook-Event-Type") != "v1.UserNotificationOpen" {
		w.WriteHeader(http.StatusOK)
		encode(w, response{
			Message: "Ok",
			Status:  true,
		})
		return
	}

	if _, err := io.Copy(hasher, r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		encode(w, response{
			Message: "Could not create crypto hash",
			Status:  false,
		})
		return
	}

	expectedHash := hex.EncodeToString(hasher.Sum(nil))

	if expectedHash != r.Header.Get("webhook-signature") {
		w.WriteHeader(http.StatusBadRequest)
		encode(w, response{
			Message: "Invalid webhook signature",
			Status:  false,
		})
		return
	}

	var request struct {
		Message string `json:"text"`
	}

	request.Message = "User opened a notification just now"

	var buf = new(bytes.Buffer)

	_ = json.NewEncoder(buf).Encode(request)

	req, err := http.NewRequest(http.MethodPost, os.Getenv("SLACK_HOOKS_URL"), buf)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encode(w, response{
			Message: "Could not send notification to Slack",
			Status:  false,
		})
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encode(w, response{
			Message: "Error while pinging Slack",
			Status:  false,
		})
		return
	}

	if resp.StatusCode > http.StatusAccepted {
		w.WriteHeader(http.StatusInternalServerError)
		encode(w, response{
			Message: "Unexpected response from Slack",
			Status:  false,
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	encode(w, response{
		Message: "Message sent to Slack successfully",
		Status:  true,
	})
}
