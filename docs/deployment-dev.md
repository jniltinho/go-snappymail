# Implantação para desenvolvimento (lab híbrido Vagrant)

Ambiente de desenvolvimento do **go-snappymail** em uma VM Ubuntu 24.04 (Vagrant + VirtualBox) com modelo **híbrido**:

- **Docker** (compose dentro da VM): apenas infraestrutura — **MariaDB** e **SnappyMail PHP** (webmail de referência).
- **Nativo** (systemd na VM): **Postfix**, **Dovecot** e os três binários Go — **go-snappymail** (buildado do repo), **go-cubemail** e **go-postfixadmin**.

## Arquitetura

```
                         Host (seu computador)
                        192.168.56.1 ─ /etc/hosts: 192.168.56.20 mail.test.local
                                 │
┌────────────────────────────────┴─────────────────────────────────────┐
│  VM Ubuntu 24.04 — 192.168.56.20 (mail.test.local)                   │
│                                                                      │
│  NATIVO (systemd)                       DOCKER (/opt/gosm-docker)    │
│  ┌────────────────────────────┐         ┌──────────────────────────┐ │
│  │ go-snappymail    :8082     │         │ gosm-mariadb             │ │
│  │ go-cubemail      :8080     │  SQL    │   127.0.0.1:3306 ──────┐ │ │
│  │ go-postfixadmin  :8081 ────┼────────►│   (db postfix)         │ │ │
│  │                            │         │                        │ │ │
│  │ Postfix  :25  ─────────────┼────────►│◄─── SQL maps ──────────┘ │ │
│  │ Dovecot  :143/:993 (SSL) ──┼────────►│                          │ │
│  │      ▲         ▲           │         │ gosm-snappymail  :8888   │ │
│  │      │ IMAP    │ LMTP/auth │  IMAP/  │   (SnappyMail PHP)       │ │
│  └──────┼─────────┴───────────┘  SMTP   │   extra_host:            │ │
│         └───────────────────────────────┤   mail.test.local ─────► │ │
│                                         └──────────────────────────┘ │
│  Maildir: /var/vmail (usuário vmail 1001:1001)                       │
│  Cert autoassinado: SANs mail.test.local, localhost, 192.168.56.20   │
└──────────────────────────────────────────────────────────────────────┘
```

Fluxo: os webmails autenticam via **IMAP no Dovecot nativo** (993/SSL) e enviam via **SMTP no Postfix nativo** (25). Postfix/Dovecot consultam domínios/mailboxes/aliases no **MariaDB em container** (publicado em `127.0.0.1:3306`), usando o mesmo schema e templates SQL de `docker/mailserver/templates/` (renderizados com `hosts = 127.0.0.1`).

## Pré-requisitos

- [Vagrant](https://www.vagrantup.com/) 2.4+ e [VirtualBox](https://www.virtualbox.org/) 7.x
- ~4 GB RAM livres
- Frontend buildado no host: `make frontend` (gera `web/dist/`, que é embedado no binário Go — o provision falha com mensagem clara se faltar)

## Subir o ambiente

```bash
make frontend          # uma vez, no host (web/dist/ é gitignored)
cd vagrant
cp .env.example .env   # opcional — só se quiser trocar senhas/versões
vagrant up
```

No `/etc/hosts` do **host**:

```
192.168.56.20  mail.test.local
```

O provision roda em ordem:

| Script | O que faz |
|--------|-----------|
| `01-base.sh` | apt: postfix, postfix-mysql, dovecot-*, mariadb-client; usuário vmail; cert autoassinado com SANs |
| `02-docker.sh` | Docker CE + compose plugin; escreve `/opt/gosm-docker/docker-compose.yml`; `docker compose up -d` (mariadb + snappymail) |
| `03-mailserver.sh` | go-postfixadmin (.deb, systemd); renderiza Postfix/Dovecot a partir de `docker/mailserver/templates` com `DB_HOST=127.0.0.1` |
| `04-go-apps.sh` | Go 1.26.x; builda go-snappymail do repo sincronizado (`/vagrant/go-snappymail`); instala go-cubemail; units systemd |
| `05-seed.sh` | Domínios/mailboxes/aliases de `docker/lab/*.txt` via CLI do postfixadmin + mysql; domínios no SnappyMail |
| `99-summary.sh` | Resumo com URLs e credenciais |

## Serviços e portas

| Serviço | Onde roda | Porta / URL |
|---------|-----------|-------------|
| go-snappymail | nativo (systemd) | http://192.168.56.20:8082 |
| go-cubemail | nativo (systemd) | http://192.168.56.20:8080 |
| go-postfixadmin | nativo (systemd) | http://192.168.56.20:8081 |
| SnappyMail PHP | Docker | http://192.168.56.20:8888 |
| MariaDB | Docker | 127.0.0.1:3306 (só dentro da VM) |
| Postfix (SMTP) | nativo | 25 |
| Dovecot (IMAP) | nativo | 143 / 993 (SSL) |

Credenciais: mailboxes com senha `Password1@` (ver `docker/lab/mailboxes.txt`), admin do postfixadmin `admin@test.local` / `Password1@`. Senha admin do SnappyMail: `docker exec gosm-snappymail cat /var/lib/snappymail/_data_/_default_/admin_password.txt`.

## Apps Go via systemd

Cada binário tem sua unit com `Restart=on-failure`:

- `/etc/systemd/system/go-snappymail.service` — binário `/opt/go-snappymail/go-snappymail`, config `/etc/go-snappymail/config.toml`
- `/etc/systemd/system/go-cubemail.service` — `/opt/go-cubemail/` (binário + config no diretório)
- `postfixadmin.service` — instalada pelo .deb em `/opt/go-postfixadmin/`

```bash
vagrant ssh
sudo systemctl status go-snappymail go-cubemail postfixadmin
sudo journalctl -u go-snappymail -f        # logs ao vivo
sudo systemctl restart go-snappymail       # após rebuild
```

Rebuild do go-snappymail depois de mudar o código (o repo é a pasta sincronizada `/vagrant/go-snappymail`):

```bash
# no host: make frontend (se mexeu no frontend/)
vagrant ssh
cd /vagrant/go-snappymail
sudo env PATH=/usr/local/go/bin:$PATH CGO_ENABLED=1 go build -o /opt/go-snappymail/go-snappymail .
sudo systemctl restart go-snappymail
```

Ou simplesmente `vagrant provision --provision-with 04-go-apps`.

## Docker compose dentro da VM

O compose vive em `/opt/gosm-docker/` (gerado pelo `02-docker.sh`), só com **mariadb** e **snappymail**. Ambos com `restart: unless-stopped` — voltam sozinhos no reboot da VM.

```bash
cd /opt/gosm-docker
docker compose ps
docker compose logs -f mariadb
docker compose restart snappymail
mysql -h127.0.0.1 -upostfix -ppostfixPassword postfix   # acesso direto ao banco
```

O container do SnappyMail tem `extra_hosts: mail.test.local:host-gateway` — assim ele alcança o Dovecot/Postfix nativos da VM.

## Seed

Reaproveitável a qualquer momento (idempotente):

```bash
vagrant provision --provision-with 05-seed
```

Fontes: `docker/lab/domains.txt`, `docker/lab/mailboxes.txt`, `docker/lab/aliases.txt` — os mesmos do lab Docker puro.

## Checklist de testes

Dentro da VM (`vagrant ssh`) ou do host, na ordem:

```bash
# 1. Units systemd ativas
systemctl is-active go-snappymail go-cubemail postfixadmin postfix dovecot

# 2. Containers up
docker compose -f /opt/gosm-docker/docker-compose.yml ps

# 3. API do go-snappymail (do host)
curl http://192.168.56.20:8082/api/v1/version

# 4. Login IMAP (na VM)
doveadm auth test 'user@test.local' 'Password1@'
# ...ou do host, IMAPS real na 993:
openssl s_client -connect 192.168.56.20:993 -quiet 2>/dev/null <<< $'a login user@test.local Password1@\na logout'

# 5. Envio SMTP (do host)
swaks --to user@test.local --from alice@test.local --server 192.168.56.20:25 \
      --auth-user alice@test.local --auth-password 'Password1@' --tls

# 6. Webmails no navegador
#    http://192.168.56.20:8082 (go-snappymail)  — user@test.local / Password1@
#    http://192.168.56.20:8888 (SnappyMail PHP) — mesmo login
```

## Troubleshooting

| Sintoma | Causa provável / correção |
|---------|---------------------------|
| Provision falha em `04-go-apps` com "web/dist/index.html não existe" | Rode `make frontend` no host e `vagrant provision` de novo |
| `go-snappymail` inativo | `journalctl -u go-snappymail -e`; config em `/etc/go-snappymail/config.toml` |
| Postfix/Dovecot não autenticam | MariaDB caiu? `docker ps`, `mariadb-admin ping -h127.0.0.1 -upostfix -ppostfixPassword` |
| SnappyMail não conecta no IMAP | O container resolve `mail.test.local` via `host-gateway`; confira `docker exec gosm-snappymail ping -c1 mail.test.local` e se o Dovecot escuta em todas as interfaces |
| Porta 3306/8888 ocupada na VM | Resto de layout antigo bare-metal — `01-base.sh` desabilita `mariadb`/`nginx` nativos; em VMs muito antigas, prefira `vagrant destroy -f && vagrant up` |
| Aviso TLS nos clientes | Esperado: cert autoassinado (SANs: `mail.test.local`, `localhost`, `192.168.56.20`) em `/etc/ssl/certs/mail.test.local.crt` |
| Seed não criou contas | `vagrant provision --provision-with 05-seed` e veja a saída; o CLI é `/opt/go-postfixadmin/postfixadmin` |

## Documentação relacionada

- [vagrant/README.md](../vagrant/README.md) — referência rápida do lab
- [docker/README.md](../docker/README.md) — lab 100% Docker (alternativa sem VM)
- [docker/LAB_ACCOUNTS.md](../docker/LAB_ACCOUNTS.md) — contas de teste
