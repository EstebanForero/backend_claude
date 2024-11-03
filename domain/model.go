package domain

// Message represents a single message in the conversation
type Message struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

// Request represents the request payload sent to the API
type Request struct {
    Model     string    `json:"model"`
    MaxTokens int       `json:"max_tokens"`
    Messages  []Message `json:"messages"`
}

// Response represents the structure of the response from the API
type Response struct {
    Content []struct {
        Text string `json:"text"`
    } `json:"content"`
}

