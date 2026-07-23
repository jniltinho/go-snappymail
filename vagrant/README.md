# go-snappymail — Ambiente de Validação (Vagrant)

VM Ubuntu 24.04 com stack completa para comparar **SnappyMail (PHP)**, **go-cubemail** e **go-snappymail** contra um servidor de email real.

## Pré-requisitos

- [Vagrant](https://www.vagrantup.com/) 2.4+
- [VirtualBox](https://www.virtualbox.org/) 7.x
- ~4 GB RAM livres

## Subir o ambiente

```bash
cd vagrant
cp .env.example .env
vagrant up
```

Adicione no `/etc/hosts` do seu computador:

```
192.168.56.20  mail.test.local
```

## Modos de instalação

### Bare-metal (padrão)

`vagrant up` instala Postfix, Dovecot, MariaDB e webmails nativamente na VM.

Serviços web: PostfixAdmin (:8081), SnappyMail (:8888), go-cubemail (:8080).

> go-snappymail (:8082) **não** é instalado no modo bare-metal — use Docker ou rode o binário manualmente na VM.

### Docker (recomendado — stack completa)

```bash
cd vagrant
vagrant provision --provision-with docker
```

Instala Docker na VM, para os serviços bare-metal e sobe a stack de [docker/README.md](../docker/README.md), incluindo **go-snappymail na porta 8082**.

Mesmas portas (`8080`, `8081`, **8082**, `8888`, `25`, `143`, `993`) — **não use os dois modos ao mesmo tempo**.

Dentro da VM:

```bash
cd /opt/go-snappymail-docker
docker compose ps
bash scripts/bootstrap.sh   # se ainda não rodou seed
```

## Serviços

| Serviço | URL | Descrição |
|---------|-----|-----------|
| **go-snappymail** | http://192.168.56.20:8082 | Este projeto (modo Docker) |
| **Go CubeMail** | http://192.168.56.20:8080 | Webmail Golang (referência arquitetural) |
| **SnappyMail (PHP)** | http://192.168.56.20:8888 | Webmail referência (UI alvo) |
| **Go-PostfixAdmin** | http://192.168.56.20:8081 | Painel admin do servidor de email |

## Credenciais de teste

| Conta | Email | Senha |
|-------|-------|-------|
| Mailbox | `user@test.local` | `Password1@` (ou `MAIL_PASS` do `.env`) |
| PostfixAdmin | `admin@test.local` | `Password1@` |
| SnappyMail admin | `admin` | `Admin1@lab` |

Lista completa: [docker/LAB_ACCOUNTS.md](../docker/LAB_ACCOUNTS.md)

## Stack instalada

### Bare-metal

- **MariaDB** — banco do PostfixAdmin (domínios, mailboxes, aliases)
- **Postfix** — MTA com virtual mailboxes via MySQL
- **Dovecot** — IMAP/LMTP/POP3 com autenticação SQL
- **Go-PostfixAdmin** v1.0.86 — gestão de domínios e contas
- **SnappyMail** v2.38.2 — nginx + PHP 8.3-FPM
- **Go CubeMail** v0.0.25 — binário Go com SPA embarcada

### Docker (provision-with docker)

Todos os serviços acima em containers, mais **go-snappymail** buildado do repositório.

## Comandos úteis

```bash
vagrant ssh                          # entrar na VM
vagrant provision --provision-with 99-summary  # reimprimir resumo
vagrant destroy -f                   # destruir VM
```

Dentro da VM (bare-metal):

```bash
doveadm auth test 'user@test.local' 'Password1@'
echo "corpo" | mail -s "Assunto teste" user@test.local
systemctl status postfix dovecot postfixadmin go-cubemail nginx
```

## Validação side-by-side

1. Acesse **SnappyMail** (8888), **go-cubemail** (8080) e **go-snappymail** (8082) com `user@test.local` / `Password1@`
2. Compare: listagem de pastas, leitura de mensagens, compose, dark mode
3. Use **PostfixAdmin** (8081) para criar mailboxes adicionais
4. Documente diferenças de UX para guiar a implementação do go-snappymail

### Screenshots com agent-browser

```bash
npm i -g agent-browser && agent-browser install

PRINTS="$(pwd)/docs/prints"
agent-browser open http://192.168.56.20:8082
agent-browser wait --load networkidle
agent-browser screenshot --full "${PRINTS}/go-snappymail-login.png"
agent-browser close
```

Prints ficam em `docs/prints/` (local, fora do git).

## Rodar go-snappymail manualmente (bare-metal)

Se preferir o binário nativo na VM em vez do container:

```bash
vagrant ssh
sudo systemctl stop go-cubemail   # opcional — libera comparação na mesma máquina
# copiar binário e config
sudo cp /vagrant/dist/go-snappymail /opt/go-snappymail/
sudo cp /vagrant/config.toml /opt/go-snappymail/
# ajustar imap.host para localhost ou mail.test.local
sudo systemctl start go-snappymail   # se unit existir
```

Ou build local na VM com Go 1.26+ e `make build-prod`.

## Estrutura

```
vagrant/
├── Vagrantfile
├── README.md
├── .env.example
├── provision/
│   ├── common.sh           # variáveis compartilhadas
│   ├── 01-base.sh          # pacotes, vmail, SSL
│   ├── 02-mariadb.sh       # banco postfix
│   ├── 03-mailserver.sh    # Postfix + Dovecot + Go-PostfixAdmin
│   ├── 04-snappymail.sh    # SnappyMail PHP
│   ├── 05-go-webmail.sh    # go-cubemail (referência Golang)
│   ├── docker.sh           # stack Docker em /opt/go-snappymail-docker
│   └── 99-summary.sh       # resumo final
└── templates/              # configs Postfix, Dovecot, nginx
```

## Documentação relacionada

- [docs/lab.md](../docs/lab.md) — visão geral do lab
- [docker/README.md](../docker/README.md) — compose e bootstrap
