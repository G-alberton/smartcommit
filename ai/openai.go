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
	openAIAPIURL = "https://api.openai.com/v1/chat/completions"
	openAIModel  = "gpt-4o-mini"
)

type OpenAIProvider struct {
	apiKey string
}

func NewOpenAIProvider() (*OpenAIProvider, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("variavel de ambiente OPENAI_API_KEY não definida")
	}
	return &OpenAIProvider{apiKey: apiKey}, nil
}

type OpenAIRequest struct {
	Model     string          `json:"model"`
	MaxTokens int             `json:"max_tokens"`
	Messages  []openAIMessage `json:"messages"`
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAiResponse struct {
	Choises []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choises"`
}

func (p *OpenAIProvider) GenerateCommitMessage(diff string) (string, error) {
	prompt := fmt.Sprintf(
		"Você é um assistente que gera mensagens de commit no padrão Conventional Commits.\n"+
			"Gere APENAS o título da mensagem, sem explicações, com no máximo 72 caracteres.\n\n"+
			"Diff:\n%s",
		diff,
	)

	reqBody := OpenAIRequest{
		Model:     openAIModel,
		MaxTokens: 100,
		Messages: []openAIMessage{
			{Role: "user", Content: prompt},
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("erro ao montar requisição: %w", err)
	}

	req, err := http.NewRequest("POST", openAIAPIURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("erro ao criar requisições: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("erro ao chamar API da OpenAI: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("erro ao ler resposta: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API retornou status %d: %s", resp.StatusCode, string(body))
	}

	var apiResp openAiResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return "", fmt.Errorf("erro ao interpretar resposta: %w", err)
	}

	if len(apiResp.Choises) == 0 {
		return "", fmt.Errorf("resposta da API sem conteudo")
	}

	return apiResp.Choises[0].Message.Content, nil
}
