# go-snappymail — Lab híbrido (Vagrant)

VM Ubuntu 24.04 com modelo **híbrido**:

- **Docker** (compose em `/opt/gosm-docker` na VM): **MariaDB** (`127.0.0.1:3306`) e **SnappyMail PHP** (`:8888`)
- **Nativo** (systemd): **Postfix** (`:25`), **Dovecot** (`:143`/`:993`) e os binários Go — **go-snappymail** (`:8082`, buildado do repo), **go-cubemail** (`:8080`), **go-postfixadmin** (`:8081`)

Guia completo (arquitetura, testes, troubleshooting): [docs/deployment-dev.md](../docs/deployment-dev.md).

## Pré-requisitos

- [Vagrant](https://www.vagrantup.com/) 2.4+ / [VirtualBox](https://www.virtualbox.org/) 7.x / ~4 GB RAM
- `make frontend` rodado no host (gera `web/dist/`, embedado no binário)

## Subir o ambiente

```bash
make frontend           # no host, uma vez
cd vagrant
cp .env.example .env    # opcional
vagrant up
```

Adicione no `/etc/hosts` do seu computador:

```
192.168.56.20  mail.test.local
```

## Serviços

| Serviço | Onde | URL / porta |
|---------|------|-------------|
| **go-snappymail** | systemd | http://192.168.56.20:8082 |
| **go-cubemail** | systemd | http://192.168.56.20:8080 |
| **go-postfixadmin** | systemd | http://192.168.56.20:8081 |
| **SnappyMail PHP** | Docker | http://192.168.56.20:8888 |
| **MariaDB** | Docker | 127.0.0.1:3306 (interno à VM) |
| **Postfix / Dovecot** | systemd | 25 / 143 / 993 |

## Credenciais de teste

| Conta | Login | Senha |
|-------|-------|-------|
| Mailbox | `user@test.local` | `Password1@` (ou `MAIL_PASS` do `.env`) |
| PostfixAdmin | `admin@test.local` | `Password1@` |
| SnappyMail admin | `admin` | `docker exec gosm-snappymail cat /var/lib/snappymail/_data_/_default_/admin_password.txt` |

Contas seed (4+ domínios): [docker/LAB_ACCOUNTS.md](../docker/LAB_ACCOUNTS.md) — mesmos arquivos `docker/lab/*.txt` do lab Docker.

## Comandos úteis

```bash
vagrant ssh                                     # entrar na VM (ou root/vagrant123)
vagrant provision --provision-with 04-go-apps   # rebuild do go-snappymail
vagrant provision --provision-with 05-seed      # re-seed (idempotente)
vagrant provision --provision-with 99-summary   # reimprimir resumo
vagrant destroy -f                              # destruir VM
```

Dentro da VM:

```bash
systemctl status go-snappymail go-cubemail postfixadmin postfix dovecot
journalctl -u go-snappymail -f
docker compose -f /opt/gosm-docker/docker-compose.yml ps
doveadm auth test 'user@test.local' 'Password1@'
curl http://localhost:8082/api/v1/version
```

## Estrutura

```
vagrant/
├── Vagrantfile
├── README.md
├── .env.example
├── provision/
│   ├── common.sh        # variáveis compartilhadas + render de templates
│   ├── 01-base.sh       # apt (postfix, dovecot, mariadb-client), vmail, cert SAN
│   ├── 02-docker.sh     # Docker CE + compose: mariadb + snappymail
│   ├── 03-mailserver.sh # Postfix/Dovecot nativos + go-postfixadmin (DB em 127.0.0.1)
│   ├── 04-go-apps.sh    # Go toolchain, build go-snappymail, go-cubemail, units systemd
│   ├── 05-seed.sh       # domínios/mailboxes/aliases de docker/lab/*.txt
│   └── 99-summary.sh    # resumo final
└── (configs Postfix/Dovecot renderizadas de ../docker/mailserver/templates)
```

## Documentação relacionada

- [docs/deployment-dev.md](../docs/deployment-dev.md) — guia completo do lab híbrido
- [docs/lab.md](../docs/lab.md) — visão geral do lab
- [docker/README.md](../docker/README.md) — alternativa 100% Docker (sem VM)
