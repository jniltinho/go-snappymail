# Screenshots — go-snappymail Lab

Capturas de tela do ambiente Vagrant (`192.168.56.20`) para comparação visual entre SnappyMail (PHP) e webmail Golang.

## Pré-requisitos

1. VM rodando: `cd vagrant && vagrant up`
2. Entrada no `/etc/hosts`: `192.168.56.20  mail.test.local`
3. [agent-browser](https://github.com/nicholasgriffintn/agent-browser) instalado:

```bash
npm i -g agent-browser && agent-browser install
```

## Gerar prints com agent-browser

```bash
PRINTS="$(pwd)/docs/prints"
BASE="http://192.168.56.20"

# SnappyMail — login
agent-browser open "${BASE}:8888"
agent-browser wait --load networkidle
agent-browser screenshot --full "${PRINTS}/01-snappymail-login.png"

# SnappyMail — inbox (após login manual ou via refs do snapshot)
agent-browser snapshot -i
agent-browser fill @e2 "user@test.local"
agent-browser fill @e3 "Password1@"
agent-browser click @e5
agent-browser wait --load networkidle
agent-browser screenshot --full "${PRINTS}/02-snappymail-inbox.png"

# Go CubeMail (referência Golang)
agent-browser open "${BASE}:8080"
agent-browser wait --load networkidle
agent-browser screenshot --full "${PRINTS}/03-go-cubemail-login.png"
# ... login e inbox

# Go SnappyMail (novo — P0 login)
agent-browser open "${BASE}:8082"
agent-browser wait --load networkidle
agent-browser screenshot --full "${PRINTS}/04-go-snappymail-login.png"
agent-browser fill @e2 "user@test.local"
agent-browser fill @e3 "Password1@"
agent-browser click @e4
agent-browser wait --load networkidle
agent-browser screenshot --full "${PRINTS}/05-go-snappymail-authenticated.png"

# Go-PostfixAdmin
agent-browser open "${BASE}:8081"
agent-browser wait --load networkidle
agent-browser screenshot --full "${PRINTS}/05-go-postfixadmin-login.png"

agent-browser close
```

## Galeria atual

| Arquivo | Descrição |
|---------|-----------|
| `01-snappymail-login.png` | SnappyMail — tela de login |
| `02-snappymail-inbox.png` | SnappyMail — inbox após login |
| `03-go-cubemail-login.png` | Go CubeMail — tela de login |
| `04-go-cubemail-inbox.png` | Go CubeMail — inbox após login |
| `05-go-postfixadmin-login.png` | Go-PostfixAdmin — login |
| `06-go-postfixadmin-dashboard.png` | Go-PostfixAdmin — dashboard |

## Credenciais

| Serviço | Login | Senha |
|---------|-------|-------|
| SnappyMail / Go CubeMail | `user@test.local` | `Password1@` |
| Go-PostfixAdmin | `admin@test.local` | `Password1@` |
| SnappyMail admin | `admin` | `Admin1@lab` |

## Rede

A VM usa rede privada VirtualBox `192.168.56.0/24`. Qualquer máquina na mesma rede host-only consegue acessar os serviços em `192.168.56.20`.

Funciona tanto no modo **bare-metal** quanto no modo **Docker** (`vagrant provision --provision-with docker`) — mesmas portas e URLs.
