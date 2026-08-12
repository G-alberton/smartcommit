package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Action int

const (
	Accept     Action = iota //0
	Edit                     //1
	Regenerate               //2
	Cancel                   //3

)

func ShowSuggestion(msg string) Action {
	fmt.Println("\nSugestão")
	fmt.Printf(" %s\n\n", msg)
	fmt.Println("[a] aceitar  [e] editar  [r] regenerar  [c] cancelar")

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("erro ao ler entrada, tente novamente")
			continue
		}

		input = strings.ToLower(strings.TrimSpace(input))

		switch input {
		case "a":
			return Accept
		case "e":
			return Edit
		case "r":
			return Regenerate
		case "c":
			return Cancel
		default:
			fmt.Println("opção inválida, digite a, e, r, ou c")
		}
	}
}
