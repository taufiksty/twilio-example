package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/twilio/twilio-go"
	twilioAPI "github.com/twilio/twilio-go/rest/api/v2010"
	"log"
	"net/http"
	"os"
	"strings"
)

var twilioClient *twilio.RestClient

// Load env
//func init() {
//	err := godotenv.Load(".env")
//	if err != nil {
//		log.Fatal("Load env error")
//	}
//}

func sendWhatsappMessage(to, message string) error {
	from := os.Getenv("TWILIO_WHATSAPP_FROM")
	if !strings.HasPrefix(to, "whatsapp:") {
		to = "whatsapp:" + to
	}

	params := &twilioAPI.CreateMessageParams{
		From: &from,
		To:   &to,
		Body: &message,
	}

	_, err := twilioClient.Api.CreateMessage(params)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}

func whatsappBotHandler(w http.ResponseWriter, r *http.Request) {
	// Parse incoming message from twilio
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unable to parse from", http.StatusBadRequest)
		return
	}

	//	Get the user's Whatsapp number and the message body
	from := r.PostFormValue("From")
	body := r.PostFormValue("Body")

	fmt.Printf("Received message from %s: %s\n", from, body)

	// Define a basic response message based on user input
	responseMessage := ""

	switch body {
	case "1":
		responseMessage = "You selected option 1. What would you like to do next? \n1. Check balance \n2. Transaction history"
	case "2":
		responseMessage = "You selected option 2. We can assist you with your order. \n1. Track order \n2. Cancel order"
	default:
		// The first interaction, providing options to the user
		responseMessage = "Welcome to our service! Please choose an option:\n1. Account Info\n2. Help with an order"
	}

	// Send the response back to the user
	err := sendWhatsappMessage(from, responseMessage)
	if err != nil {
		log.Printf("failed to send message: %s\n", err)
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		return
	}

	// Acknowledge Twilio with a 200 OK response
	w.WriteHeader(http.StatusOK)
}

func main() {
	twilioAccountSID := os.Getenv("TWILIO_ACCOUNT_SID")
	twilioAuthToken := os.Getenv("TWILIO_AUTH_TOKEN")

	twilioClient = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: twilioAccountSID,
		Password: twilioAuthToken,
	})

	// Create a new router
	router := mux.NewRouter()

	// Define a route for Twilio webhook (to receive incoming WhatsApp messages)
	router.HandleFunc("/whatsapp-webhook", whatsappBotHandler).Methods("POST")

	// Start the server
	port := os.Getenv("PORT")
	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
