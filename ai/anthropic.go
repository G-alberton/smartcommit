package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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
	if apiKey == "" {
		return nil, fmt.Errorf("Variavel de ambiente ANTHROPIC_API_KEY não definida")
	}
	return &AnthropicProvider{apiKey: apiKey}, nil
}

type anthropicRequest struct {
	Model     string             `json:"model"`
	MaxTokens int                `json:"max_tokens"`
	Messages  []anthropicMessage `json:"messages"`
}

type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type anthropicResponse struct {
	Content []struct {
		Type string `json:"Type"`
		Text string `json:"text"`
	} `json:"content"`
}

func (p *AnthropicProvider) GenerateCommitMessage(diff string) (string, error) {
	prompt := fmt.Sprintf(
		"Você é um assistente que gera mensagens de commit no padrão Conventional Commits. \n"+
			"Gere APENAS o Título da mensagem, sem explicação, com no máximo 72 caracteres. \n\n"+
			"Diff:\n%s",
		diff,
	)

	reqBody := anthropicRequest{
		Model:     anthropicModel,
		MaxTokens: 100,
		Messages: []anthropicMessage{
			{Role: "user", Content: prompt},
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("erro ao montar requisição: %w", err)
	}

	req, err := http.NewRequest("POST", anthropicAPIURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("erro ao criar requisições: %w", err)
	}

	req.Header.Set("x-api-key", p.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("erro ao chamar API da Anthoropic: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("erro ao ler resposta: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API retornou status %d: %s", resp.StatusCode, string(body))
	}

	var apiResp anthropicResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return "", fmt.Errorf("erro ao interpretar resposta: %w", err)
	}

	if len(apiResp.Content) == 0 {
		return "", fmt.Errorf("resposta da APi sem conteudo")
	}

	return apiResp.Content[0].Text, nil
}
