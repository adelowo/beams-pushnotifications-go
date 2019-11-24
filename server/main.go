package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	pushnotifications "github.com/pusher/push-notifications-go"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := flag.Int64("http.port", 3000, "Port to run HTTP server on")

	flag.Parse()

	beamsClient, err := pushnotifications.New(os.Getenv("PUSHER_BEAMS_INSTANCE_ID"), os.Getenv("PUSHER_BEAMS_SECRET_KEY"))
	if err != nil {
		log.Fatalf("Could not set up Push Notifications client... %v", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/push", createPushNotificationHandler(beamsClient))
	mux.HandleFunc("/auth", authenticateUser(beamsClient))
	mux.HandleFunc("/slack", handleWebhook)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), mux); err != nil {
		log.Fatal(err)
	}
}

var currentUser = ""

func authenticateUser(client pushnotifications.PushNotifications) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userIDinQueryParam := r.URL.Query().Get("user_id")

		beamsToken, err := client.GenerateToken(userIDinQueryParam)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		currentUser = userIDinQueryParam

		beamsTokenJson, err := json.Marshal(beamsToken)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(beamsTokenJson)
	}
}

func createPushNotificationHandler(client pushnotifications.PushNotifications) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var data map[string]interface{}

		type response struct {
			Status  bool   `json:"status"`
			Message string `json:"message"`
		}

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			encode(w, response{
				Status:  false,
				Message: "Invalid bad request",
			})
			return
		}

		publishRequest := map[string]interface{}{
			"apns": map[string]interface{}{
				"aps": map[string]interface{}{
					"alert": data,
				},
			},
			"fcm": map[string]interface{}{
				"notification": data,
			},
		}

		_, err := client.PublishToUsers([]string{currentUser}, publishRequest)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			encode(w, response{
				Status:  false,
				Message: "Could not send push notification",
			})
			return
		}

		w.WriteHeader(http.StatusOK)
		encode(w, response{
			Status:  true,
			Message: "Push notification sent successfully",
		})
	}
}

var encode = func(w http.ResponseWriter, v interface{}) {
	_ = json.NewEncoder(w).Encode(v)
}
