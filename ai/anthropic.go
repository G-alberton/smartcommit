package ai

import (
	"net/http"
	"encoding/json"
	"bytes"
	"io"
	"os"
	"fmt"
)

const (
	anthropicAPIURL = "https://api.anthropic.com/v1/messages"
	anthropicModel  = "claude-sonnet-4-5-20250929"
)

type AnthropicProvider struct {
	apiKey string
}

func NewAnthropicProvider() (*AnthropicProvider, error) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey  == "" {
		return nil, fmt.Errorf("Variavel de ambiente ANTHROPIC_API_KEY não definida")
	}
	return &AnthropicProvider{apiKey: apiKey}, nil
}

type anthropicRequest struct {
	Model string `json:"model"`
	MaxTokens int `json:"max_tokens"`
	Messages []anthropicMessage `json:"messages"`
}

type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type anthropicResponse

