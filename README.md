# smartcommit
CLI em GO que gera mensagem de commit no padrão Conventional Commits, usando IA para analisar o 'diff' das mudanças staged.

## Motivação
Garantir que os Commits sejam e estejam feito de maneira correta sem negligenciar a organização dos nossos codigos.
Vc faz o `git add` normalmente e a ferramenta le o diff, manda para a IA e sugere uma mensagem, voce decide se vai aceitar, ou não.

## Como funciona

1. Você roda `git add` normalmente.
2. `smartcommit` lê o diff staged (`git diff --cached`).
3. O diff é enviado como prompt para a IA configurada.
4. A IA devolve uma sugestão de mensagem no padrão Conventional Commits.
5. Você escolhe: **[a]** aceitar, **[e]** editar, **[r]** regenerar ou **[c]** cancelar.
6. Se aceito, o commit é criado automaticamente (`git commit -m "..."`).

# Por que GO?
- **Estudos**: Quis fazer o projeto para que seja possivel estudar enquanto produzia algo legal, tentando usar o minimo de IA
- **`os/exec` nativo**: Permite chamar o `git` diretamente do sistema, sem dependecia
- **`net/http` + `encoding/json` stdlib**: As chamadas de API são da biblioteca padrão
-**Cross**: Funcionamento em diferentes sistema operacionais

## Instalação
 `git clone https://github.com/G-alberton/smartcommit`
 `cd smartcommit`
 `go build -o smartcommit`

 ## Configuração
 `export ANTHROPIC_API_KEY="sua_chave"`
 ou
 `export OPENAI_API_KEY="sua_chave"`