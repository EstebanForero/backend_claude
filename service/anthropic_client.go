package service

import (
    "bytes"
    "encoding/json"
    "io/ioutil"
    "net/http"
    "backend_claude/domain"
)

// AnthropicClient represents the client for Anthropic API
type AnthropicClient struct {
    apiKey string
}

// NewAnthropicClient creates a new instance of AnthropicClient
func NewAnthropicClient(apiKey string) *AnthropicClient {
    return &AnthropicClient{
        apiKey: apiKey,
    }
}

// SendMessage sends a message to the Anthropic API and returns the response
func (c *AnthropicClient) SendMessage(content string) (string, error) {
    url := "https://api.anthropic.com/v1/messages"

    // Prepare the request payload
    message := domain.Request{
        Model:     "claude-3-5-sonnet-20241022",
        MaxTokens: 1024,
        Messages: []domain.Message{
            {
                Role:    "user",
                Content: content,
            },
        },
    }

    jsonData, err := json.Marshal(message)
    if err != nil {
        return "", err
    }

    // Create a new HTTP request
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return "", err
    }

    // Set the necessary headers
    req.Header.Set("x-api-key", c.apiKey)
    req.Header.Set("anthropic-version", "2023-06-01")
    req.Header.Set("Content-Type", "application/json")

    // Send the request
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    // Read the response
    bodyBytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    return string(bodyBytes), nil
}

