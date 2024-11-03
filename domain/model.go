package domain

type Message struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

type Request struct {
    Model     string    `json:"model"`
    MaxTokens int       `json:"max_tokens"`
    Messages  []Message `json:"messages"`
}

type Response struct {
    Content []struct {
        Text string `json:"text"`
    } `json:"content"`
}

