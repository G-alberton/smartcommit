package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"smartcommit/ai"
	"smartcommit/cli"
	"smartcommit/git"
)

func main() {
	providerName := flag.String("provider", "anthropic", "provider de IA a usar (anthropic, openai, deepseek)")
	flag.Parse()

	diff, err := git.GetStagedDiff()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	provider, err := ai.NewProvider(*providerName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for {
		fmt.Println("Analisando mudanças...")

		msg, err := provider.GenerateCommitMessage(diff)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		action := cli.ShowSuggestion(msg)

		switch action {
		case cli.Accept:
			commit(msg)
			return
		case cli.Edit:
			fmt.Print("Digite a nova mensagem")
			reader := bufio.NewReader(os.Stdin)
			newMsg, _ := reader.ReadString('\n')
			msg = strings.TrimSpace(newMsg)
			commit(msg)
			return

		case cli.Regenerate:
			continue

		case cli.Cancel:
			fmt.Println("Cancelado")
			os.Exit(0)
		}
	}
}

func commit(msg string) {
	cmd := exec.Command("git", "commit", "-m", msg)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		os.Exit(1)
	}
	fmt.Printf("Comit criado: %s\n", msg)
}
