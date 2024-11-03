package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "backend_claude/service"
)

type requestBody struct {
    Message string `json:"message"`
}

type responseBody struct {
    Response string `json:"response"`
}

func main() {
    // Retrieve the API key from an environment variable
    apiKey := os.Getenv("ANTHROPIC_API_KEY")
    if apiKey == "" {
        log.Fatal("ANTHROPIC_API_KEY environment variable not set")
    }

    // Create a new Anthropic client
    client := service.NewAnthropicClient(apiKey)

    // Set up the HTTP server
    http.HandleFunc("/ask-claude", func(w http.ResponseWriter, r *http.Request) {
        // Ensure the request method is POST
        if r.Method != http.MethodPost {
            http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
            return
        }

        // Parse the JSON request body
        var reqBody requestBody
        err := json.NewDecoder(r.Body).Decode(&reqBody)
        if err != nil {
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }

        // Send the message to the Anthropic API
        response, err := client.SendMessage(reqBody.Message)
        if err != nil {
            http.Error(w, "Failed to get response from Anthropic API", http.StatusInternalServerError)
            return
        }

        // Return the response as JSON
        resBody := responseBody{Response: response}
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(resBody)
    })

    fmt.Println("Server is running on http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

