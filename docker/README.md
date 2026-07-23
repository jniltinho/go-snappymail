# Stack Docker — go-snappymail Lab

Ambiente containerizado equivalente ao provisionamento bare-metal da VM Vagrant.

Documentação completa: [docs/lab.md](../docs/lab.md) · contas de teste: [LAB_ACCOUNTS.md](LAB_ACCOUNTS.md)

## Serviços

| Container | Porta | Imagem |
|-----------|-------|--------|
| `gosm-mariadb` | 3306 (interno) | mariadb:10.11 |
| `gosm-postfixadmin` | 8081 | build local (`postfixadmin/`) |
| `gosm-mailserver` | 25, 143, 993 | build local (`mailserver/`) |
| `gosm-snappymail` | 8888 | djmaze/snappymail:latest |
| `gosm-go-cubemail` | 8080 | build local (`go-cubemail/`) |
| `gosm-go-snappymail` | **8082** | build local (repo raiz) |

## Uso na VM Vagrant

```bash
# Modo Docker (para serviços bare-metal e sobe containers)
cd vagrant
vagrant provision --provision-with docker

# Ou manualmente dentro da VM
cd /opt/go-snappymail-docker
docker compose up -d
bash scripts/bootstrap.sh
```

## Uso no host (sem Vagrant)

Requer Docker + Docker Compose v2:

```bash
cd docker
docker compose up -d --build
bash scripts/bootstrap.sh
```

Adicione no `/etc/hosts`: `127.0.0.1 mail.test.local` (host) ou `192.168.56.20 mail.test.local` (VM).

Seed all lab domains and mailboxes:

```bash
bash scripts/bootstrap.sh
# or re-seed accounts only:
bash scripts/seed-lab.sh
```

## Credentials

See `.env` and **[LAB_ACCOUNTS.md](LAB_ACCOUNTS.md)** for all test domains and mailboxes.

Default password for every lab mailbox: **`Password1@`**

- PostfixAdmin superadmin: `admin@test.local` / `Password1@`

## Prints com agent-browser

Mesmas URLs na rede `192.168.56.20`:

```bash
agent-browser open http://192.168.56.20:8888
agent-browser screenshot --full ../docs/prints/snappymail-docker-login.png

agent-browser open http://192.168.56.20:8082
agent-browser screenshot --full ../docs/prints/go-snappymail-docker-login.png
```

## Comandos úteis

```bash
docker compose ps
docker compose logs -f mailserver
docker compose exec mailserver doveadm auth test 'user@test.local' 'Password1@'
docker compose down -v   # reset completo (apaga volumes)
```

## Arquitetura

```
┌─────────────┐  ┌──────────────┐  ┌──────────────┐  ┌─────────────┐
│ SnappyMail  │  │  Go CubeMail │  │ go-snappymail│  │ PostfixAdmin│
│   :8888     │  │    :8080     │  │    :8082     │  │    :8081    │
└──────┬──────┘  └──────┬───────┘  └──────┬───────┘  └──────┬──────┘
       │                │                  │
       └────────────────┼──────────────────┘
                        │  rede gosm (bridge)
              ┌─────────┴─────────┐
              │    mailserver     │
              │  Postfix+Dovecot  │
              │   :25/:143/:993   │
              └─────────┬─────────┘
                        │
              ┌─────────┴─────────┐
              │     MariaDB       │
              └───────────────────┘
```

## Modos na VM

| Modo | Comando | Descrição |
|------|---------|-----------|
| **bare-metal** | `vagrant up` (default) | Postfix/Dovecot/nginx nativos |
| **docker** | `vagrant provision --provision-with docker` | Tudo em containers |

Os dois modos usam as mesmas portas — não rode os dois simultaneamente.
