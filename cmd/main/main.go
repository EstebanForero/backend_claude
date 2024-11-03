package main

import (
	"fmt"
	"log"
	"os"
    "backend_claude/service"
)

func main() {
    // Retrieve the API key from an environment variable
    apiKey := os.Getenv("ANTHROPIC_API_KEY")
    if apiKey == "" {
        log.Fatal("ANTHROPIC_API_KEY environment variable not set")
    }

    // Create a new Anthropic client
    client := service.NewAnthropicClient(apiKey)

    // Send a message to the Claude model
    response, err := client.SendMessage("Hello, world")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Response from Claude:", response)
}

