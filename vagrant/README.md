# go-snappymail — Ambiente de Validação (Vagrant)

VM Ubuntu 24.04 com stack completa para comparar **SnappyMail (PHP)** vs **webmail Golang** contra um servidor de email real.

## Pré-requisitos

- [Vagrant](https://www.vagrantup.com/) 2.4+
- [VirtualBox](https://www.virtualbox.org/) 7.x
- ~4 GB RAM livres

## Subir o ambiente

```bash
cd vagrant
vagrant up
```

Adicione no `/etc/hosts` do seu computador:

```
192.168.56.20  mail.test.local
```

## Modos de instalação

### Bare-metal (padrão)

`vagrant up` instala Postfix, Dovecot, MariaDB e webmails nativamente na VM.

### Docker (containers)

```bash
cd vagrant
vagrant provision --provision-with docker
```

Instala Docker na VM, para os serviços bare-metal e sobe a stack de [docker/README.md](../docker/README.md). Mesmas portas (`8080`, `8081`, `8888`, `25`, `143`, `993`) — **não use os dois modos ao mesmo tempo**.

## Serviços

| Serviço | URL | Descrição |
|---------|-----|-----------|
| **Go-PostfixAdmin** | http://192.168.56.20:8081 | Painel admin do servidor de email |
| **SnappyMail (PHP)** | http://192.168.56.20:8888 | Webmail referência (UI alvo) |
| **Go CubeMail** | http://192.168.56.20:8080 | Webmail Golang (referência arquitetural) |

> **Nota:** `go-snappymail` ainda não tem código implementado. O **go-cubemail** ocupa a porta 8080 como referência Golang até o P0 do go-snappymail estar pronto.

## Credenciais de teste

| Conta | Email | Senha |
|-------|-------|-------|
| Mailbox | `user@test.local` | `Password1@` |
| PostfixAdmin | `admin@test.local` | `Password1@` |
| SnappyMail admin | `admin` | `Admin1@lab` |

## Stack instalada

- **MariaDB** — banco do PostfixAdmin (domínios, mailboxes, aliases)
- **Postfix** — MTA com virtual mailboxes via MySQL
- **Dovecot** — IMAP/LMTP/POP3 com autenticação SQL
- **Go-PostfixAdmin** v1.0.86 — gestão de domínios e contas
- **SnappyMail** v2.38.2 — nginx + PHP 8.3-FPM
- **Go CubeMail** v0.0.25 — binário Go com SPA embarcada

## Comandos úteis

```bash
vagrant ssh                          # entrar na VM
vagrant provision --provision-with 99-summary  # reimprimir resumo
vagrant destroy -f                   # destruir VM
```

Dentro da VM:

```bash
# Testar autenticação IMAP
doveadm auth test 'user@test.local' 'Password1@'

# Enviar email de teste
echo "corpo" | mail -s "Assunto teste" user@test.local

# Status dos serviços
systemctl status postfix dovecot postfixadmin go-cubemail nginx
```

## Validação side-by-side

1. Acesse **SnappyMail** (8888) e **Go CubeMail** (8080) com `user@test.local` / `Password1@`
2. Compare: listagem de pastas, leitura de mensagens, compose, dark mode
3. Use **Go-PostfixAdmin** (8081) para criar mailboxes adicionais
4. Documente diferenças de UX para guiar a implementação do go-snappymail

### Screenshots com agent-browser

A VM fica na rede `192.168.56.0/24` — acessível do host em `192.168.56.20`. Use [agent-browser](https://github.com/nicholasgriffintn/agent-browser) para capturar prints em `docs/prints/`:

```bash
npm i -g agent-browser && agent-browser install

PRINTS="$(pwd)/docs/prints"
agent-browser open http://192.168.56.20:8888
agent-browser wait --load networkidle
agent-browser screenshot --full "${PRINTS}/snappymail-login.png"
agent-browser close
```

Veja o guia completo em [docs/prints/README.md](../docs/prints/README.md).

## Estrutura

```
vagrant/
├── Vagrantfile
├── README.md
├── provision/
│   ├── common.sh           # variáveis compartilhadas
│   ├── 01-base.sh          # pacotes, vmail, SSL
│   ├── 02-mariadb.sh       # banco postfix
│   ├── 03-mailserver.sh    # Postfix + Dovecot + Go-PostfixAdmin
│   ├── 04-snappymail.sh    # SnappyMail PHP
│   ├── 05-go-webmail.sh    # go-cubemail (referência Golang)
│   └── 99-summary.sh       # resumo final
└── templates/              # configs Postfix, Dovecot, nginx
```

## Substituir go-cubemail por go-snappymail

Quando o binário `go-snappymail` estiver disponível:

```bash
# Na VM ou via synced folder
systemctl stop go-cubemail
cp /vagrant/go-snappymail/bin/go-snappymail /opt/go-snappymail/
# ajustar config.toml e systemd unit
systemctl start go-snappymail
```
