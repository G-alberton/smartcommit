package ai

import "fmt"

func NewProvider(name string) (Provider, error) {
	switch name {
	case "anthropic":
		return NewAnthropicProvider()
	case "openai":
		return NewOpenAIProvider()
	case "deepSeek":
		return NewDeepSeekProvider()
	default:
		return nil, fmt.Errorf("provider desconhecido: %s", name)
	}
}
