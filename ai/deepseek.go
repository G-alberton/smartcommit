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
	deepSeekAPIURL = "https://api.deepseek.com/chat/completions"
	deepSeekModel  = "deepseek-chat"
)

type DeepSeekProvider struct {
	apiKey string
}

func NewDeepSeekProvider() (*DeepSeekProvider, error) {
	apiKey := os.Getenv("DEPSEEK_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("variavel de ambiente DEEPSEEK_API_KEY não definida")
	}
	return &DeepSeekProvider{apiKey: apiKey}, nil
}

type deepSeekRequest struct {
	Model     string            `json:"model"`
	MaxTokens int               `json:"max_tokens"`
	Messages  []deepSeekMessage `json:"messages"`
}

type deepSeekMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type deepseekResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func (p *DeepSeekProvider) GenerateCommitMessage(diff string) (string, error) {
	prompt := fmt.Sprintf(
		"Você é um assistente que gera mensagens de commit no padrão Conventional Commits.\n"+
			"Gere APENAS o título da mensagem, sem explicações, com no máximo 72 caracteres.\n\n"+
			"Diff:\n%s",
		diff,
	)

	reqBody := deepSeekRequest{
		Model:     deepSeekModel,
		MaxTokens: 100,
		Messages: []deepSeekMessage{
			{Role: "user", Content: prompt},
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("erro ao montar requisição: %w", err)
	}

	req, err := http.NewRequest("POST", deepSeekAPIURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("erro ao criar requisição: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("erro ao chamar API da DeepSeek: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("erro ao ler resposta: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API retornou status %d: %s", resp.StatusCode, string(body))
	}

	var apiResp deepseekResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return "", fmt.Errorf("erro ao interpretar resposta: %w", err)
	}

	if len(apiResp.Choices) == 0 {
		return "", fmt.Errorf("resposta da API sem conteúdo")
	}

	return apiResp.Choices[0].Message.Content, nil
}
