package service

import (
	"backend_claude/domain"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

const systemMessage = "You are ISA, an AI assistant for Belcorp. You can help users with the following:\n\n" +
    "1. **Product Information**: Search and retrieve information about Belcorp's products.\n" +
    "2. **Recommendations**: Provide the top recommended products.\n" +
    "3. **Order Processing**: Assist with adding items to an order, confirming the order, and completing the purchase.\n\n" +
    "When a user wants to place an order, ask them for product names and quantities, then add them to the order and provide a summary before proceeding with the purchase."


type AnthropicClient struct {
    apiKey string
}

func NewAnthropicClient(apiKey string) *AnthropicClient {
    return &AnthropicClient{
        apiKey: apiKey,
    }
}

func (c *AnthropicClient) SendMessage(messages []domain.Message) (string, error) {
    url := "https://api.anthropic.com/v1/messages"

    messages = append(messages, domain.Message{
        Role: "system",
        Content: systemMessage,
    })

    message := domain.Request{
        Model:     "claude-3-5-sonnet-20241022",
        MaxTokens: 1024,
        Messages:  messages,
    }

    jsonData, err := json.Marshal(message)
    if err != nil {
        return "", err
    }

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return "", err
    }

    req.Header.Set("x-api-key", c.apiKey)
    req.Header.Set("anthropic-version", "2023-06-01")
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    bodyBytes, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    var apiResponse domain.Response
    if err := json.Unmarshal(bodyBytes, &apiResponse); err != nil {
        return "", err
    }

    if len(apiResponse.Content) > 0 {
        return apiResponse.Content[0].Text, nil
    }

    return "", errors.New("no text content found in response")
}

