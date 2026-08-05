package ai

type Provider interface {
	GenerateCommitMessage(diff string) (string, error)
}
