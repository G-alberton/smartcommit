# smartcommit

CLI em Go que gera mensagens de commit no padrão Conventional Commits, usando IA para analisar o `diff` das mudanças staged.

## Motivação

Garantir que os commits sejam feitos de maneira correta, sem negligenciar a organização do código. Você faz o `git add` normalmente, e a ferramenta lê o diff, manda para a IA e sugere uma mensagem — você decide se aceita ou não.

## Como funciona

1. Você roda `git add` normalmente.
2. `smartcommit` lê o diff staged (`git diff --cached`).
3. O diff é enviado como prompt para a IA configurada.
4. A IA devolve uma sugestão de mensagem no padrão Conventional Commits.
5. Você escolhe: **[a]** aceitar, **[e]** editar, **[r]** regenerar ou **[c]** cancelar.
6. Se aceito, o commit é criado automaticamente (`git commit -m "..."`).

## Providers suportados

| Provider  | Flag `--provider`    | Variável de ambiente |
|-----------|-----------------------|------------------------|
| Anthropic | `anthropic` (padrão) | `ANTHROPIC_API_KEY`   |
| OpenAI    | `openai`              | `OPENAI_API_KEY`      |
| DeepSeek  | `deepseek`            | `DEEPSEEK_API_KEY`    |

Você só precisa configurar a chave do provider que for usar.

## Por que Go?

- **Estudo**: quis fazer esse projeto pra estudar Go e arquitetura de software enquanto produzia algo útil, tentando usar o mínimo de IA na escrita do próprio código.
- **`os/exec` nativo**: permite chamar o `git` diretamente do sistema, sem dependência externa.
- **`net/http` + `encoding/json` na stdlib**: as chamadas de API usam só a biblioteca padrão, sem bibliotecas de terceiros.
- **Cross-platform**: compila e roda em diferentes sistemas operacionais sem alterações.

## Instalação

Requer o [Go](https://go.dev/dl/) instalado (versão 1.21+).

### Opção 1 — `go install` (mais simples)

```bash
go install github.com/G-alberton/smartcommit@latest
```

Isso baixa, compila e instala o `smartcommit` em `$GOPATH/bin` (ou `$HOME/go/bin`) numa linha só. Garanta que esse diretório esteja no seu `PATH` pra poder rodar `smartcommit` de qualquer lugar.

### Opção 2 — clonar e compilar localmente

Útil se você quer mexer no código ou não quer instalar globalmente.

```bash
git clone https://github.com/G-alberton/smartcommit
cd smartcommit
go build
```

Isso gera o executável (`smartcommit.exe` no Windows, `smartcommit` no Linux/Mac) na pasta atual.

## Configuração

Defina a variável de ambiente correspondente ao provider que você quer usar.

**Linux/macOS (bash/zsh):**

```bash
export ANTHROPIC_API_KEY="sua_chave"
```

**Windows (PowerShell):**

```powershell
$env:ANTHROPIC_API_KEY = "sua_chave"
```

## Uso

```bash
git add .
smartcommit
```

Por padrão, usa a Anthropic. Para escolher outro provider:

```bash
smartcommit --provider openai
smartcommit --provider deepseek
```

### Exemplo

```
$ git add .
$ smartcommit
Analisando mudanças...

Sugestão:
  feat(auth): adiciona validação de token JWT no middleware

[a] aceitar  [e] editar  [r] regenerar  [c] cancelar
> a

✔ Commit criado: feat(auth): adiciona validação de token JWT no middleware
```

## Limitações conhecidas

- Sem arquivo de configuração (a escolha de provider é só via flag).
- Mensagem de commit é só o título (sem corpo multi-linha).
- Diffs muito grandes podem ultrapassar o limite de tokens da API.
