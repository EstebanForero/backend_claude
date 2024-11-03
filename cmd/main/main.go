package main

import (
	"backend_claude/service"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"backend_claude/domain"
	"github.com/rs/cors"
)

type requestBody struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}

type responseBody struct {
	Response string `json:"response"`
}

func main() {

	var (
		conversationHistory = make(map[string][]domain.Message)
		historyMutex        = &sync.RWMutex{}
	)

	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		log.Fatal("ANTHROPIC_API_KEY environment variable not set")
	}

	client := service.NewAnthropicClient(apiKey)

	// Define the handler function
	http.HandleFunc("/ask-claude", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.WriteHeader(http.StatusOK)
			return
		}

		log.Println("Received a request")

		if r.Method != http.MethodPost {
			http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
			log.Println("Invalid method")
			return
		}

		// Parse the JSON request body
		var reqBody requestBody
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			log.Println("Failed to decode request body:", err)
			return
		}
		log.Println("Request body decoded successfully:", reqBody)

		if reqBody.UserID == "" {
			http.Error(w, "User ID is required", http.StatusBadRequest)
			log.Println("User ID not provided")
			return
		}

		// Lock the history and add the user's message
		historyMutex.Lock()
		userHistory := conversationHistory[reqBody.UserID]

		userHistory = append(userHistory, domain.Message{
			Role:    "user",
			Content: reqBody.Message,
		})

		if len(userHistory) > 10 {
			userHistory = userHistory[len(userHistory)-10:]
		}

		conversationHistory[reqBody.UserID] = userHistory
		historyMutex.Unlock()

		// Send the message to the API and get a response
		response, err := client.SendMessage(userHistory)

		if err != nil {
			http.Error(w, "Failed to get response from Anthropic API", http.StatusInternalServerError)
			log.Println("Error from Anthropic API:", err)
			return
		}
		log.Println("Response from Anthropic API:", response)

		// Update conversation history with the assistant's response
		historyMutex.Lock()
		conversationHistory[reqBody.UserID] = append(conversationHistory[reqBody.UserID], domain.Message{
			Role:    "assistant",
			Content: response,
		})
		if len(conversationHistory[reqBody.UserID]) > 10 {
			conversationHistory[reqBody.UserID] = conversationHistory[reqBody.UserID][len(conversationHistory[reqBody.UserID])-10:]
		}
		historyMutex.Unlock()

		// Send the response back
		resBody := responseBody{Response: response}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resBody)
	})

	// Wrap the default server with CORS
	c := cors.Default()
	handler := c.Handler(http.DefaultServeMux)

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

