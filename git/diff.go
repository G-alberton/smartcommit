package git

import (
	"fmt"
	"os/exec"
	"strings"
)

func GetStagedDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--cached") //aqui a gente monta o comando para "rodar" no termina e pegar a diff do codigo feito

	output, err := cmd.Output() //bytes de saida do comando, aqui que roda o comando anterior
	if err != nil {
		return "", fmt.Errorf("erro ao executar git diff: %w", err)
	}

	diff := strings.TrimSpace(string(output))
	if diff == "" {
		return "", fmt.Errorf("nenhuma mudança staged encontrada - rode 'git add' primeiro")
	}

	return diff, nil

}
