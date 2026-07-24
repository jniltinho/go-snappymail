# Ambiente de desenvolvimento — "Zimbra em Go"

Como subir todo o ambiente para desenvolver o ecossistema (webmail + painel admin),
com o **Zimbra de referência já instalado** (via snapshot, sem reinstalar do zero).

---

## Visão geral

O ecossistema recria o Zimbra em Go + Vue 3, em peças independentes:

| Peça | Papel | Porta | Repo |
|---|---|---|---|
| **go-snappymail** | Webmail (clone do Zimbra Web Classic) | 8082 | este |
| **painel-admin** | Admin (clone do ZimbraAdmin) | 7071 (http/https) | go-postfixadmin |
| **mailserver** | Postfix + Dovecot | 25/143/993/587 | Docker / VM |
| **banco** | MariaDB ou PostgreSQL (GORM) | 3306 / 5432 | Docker / VM |
| **zimbra (ref)** | Zimbra FOSS instalado (referência) | 443 / 7071 | vagrant/zimbra |

Cada binário embute sua própria SPA Vue (`go:embed`) e não depende dos outros em runtime;
compartilham apenas o servidor de mail e o banco.

---

## 1. Zimbra de referência — subir JÁ INSTALADO (snapshot)

A VM `vagrant/zimbra` tem o **Zimbra FOSS 10.1.17.p2** instalado. A instalação do zero
leva 20-40 min; para desenvolver, use o **snapshot** `zimbra-installed` (sobe em ~1-2 min).

### Primeira vez (criar a VM e o snapshot)

```bash
cd vagrant/zimbra
vagrant up                              # instala do zero (20-40 min, baixa 236MB)
vagrant snapshot save zimbra-installed  # congela o estado "tudo instalado"
```

### Dia a dia (restaurar o snapshot)

```bash
cd vagrant/zimbra
vagrant snapshot restore zimbra-installed   # volta ao estado instalado em 1-2 min
```

Outros comandos úteis:

```bash
vagrant snapshot list                       # lista snapshots
vagrant halt                                # desliga (preserva estado)
vagrant snapshot restore zimbra-installed   # descarta alterações e volta ao snapshot
vagrant ssh -c "sudo su - zimbra -c 'zmcontrol status'"   # status dos serviços
```

No `/etc/hosts` do **host**:

```
192.168.56.30  mail.zimbra.test
```

| Acesso | URL | Credenciais |
|---|---|---|
| Webmail (referência) | https://192.168.56.30/ | `nilton@linuxpro.com.br` / `Password1@` |
| Admin console (referência) | https://mail.zimbra.test:7071/ | `admin@zimbra.test` / `Password1@` |

Estado semeado: domínios `zimbra.test`, `linuxpro.com.br` (34 contas), `criarenet.com`
(20 contas); a caixa `nilton@linuxpro.com.br` tem ~500 mensagens de teste (com anexos PDF/HTML)
para exercícios de UI. Administração via `zmprov`: ver [`zimbra-vagrant-foss.md`](zimbra-vagrant-foss.md).

---

## 2. go-snappymail (webmail)

```bash
# build + run (SQLite local, config.toml)
make run                 # serve http://localhost:8082

# frontend com HMR (Vite :5173, proxy → :8082)
make frontend-dev
```

### Apontar o webmail para a VM Zimbra (mesma caixa da referência)

Edite o `config.toml` (ou use env `GOSM_*`):

```toml
[imap]
host = "192.168.56.30"
port = 993
tls  = true
tls_server_name = "mail.zimbra.test"
insecure_skip_verify = true   # lab: cert self-signed

[smtp]
host = "192.168.56.30"
port = 587
starttls = true
insecure_skip_verify = true
```

Assim o go-snappymail e o Zimbra de referência mostram **a mesma caixa** — ideal para
comparar o clone lado a lado. Login: `nilton@linuxpro.com.br` / `Password1@`.

### Lab Docker alternativo

Sem a VM, use o mailserver dockerizado (`docker/docker-compose.yml`): IMAP `:143/:993`,
app `:8082`, usuário `user@test.local` / `Password1@`. Ver [`lab.md`](lab.md).

---

## 3. painel-admin (fase atual — ZimbraAdmin em Go/Vue)

Backend em **go-postfixadmin** (GORM → MariaDB/PostgreSQL, JWT/RBAC), frontend Vue 3 com
layout idêntico ao ZimbraAdmin, servido em `:7071` (http ou https).

```bash
cd ../../go-postfixadmin           # repo irmão
make build                         # Vue admin + binário Go embedded
./bin/postfixadmin --generate-config
./bin/postfixadmin migrate && ./bin/postfixadmin migrate rbac
./bin/postfixadmin server          # painel admin :7071
```

Banco de dev (Docker): MariaDB `:3306` do `docker/docker-compose.yml` do go-snappymail,
ou o compose do próprio go-postfixadmin. Detalhes e roadmap: proposta OpenSpec
`admin-panel-zimbra` em `go-postfixadmin/openspec/changes/`.

---

## Fluxo recomendado (tudo de uma vez)

```bash
# 1. Zimbra de referência (snapshot)
(cd vagrant/zimbra && vagrant snapshot restore zimbra-installed)

# 2. Webmail apontando para a VM
make run                    # :8082  (config → 192.168.56.30)

# 3. Painel admin
(cd ../go-postfixadmin && ./bin/postfixadmin server)   # :7071

# 4. Comparar clones com a referência
#    webmail:  http://localhost:8082   vs  https://192.168.56.30/
#    admin:    http://localhost:7071   vs  https://mail.zimbra.test:7071/
```

---

## Screenshots (convenção)

Ao validar UI ou comparar clones com a referência, **sempre** salve os prints em
`docs/prints/` (subpasta por contexto), nunca em `/tmp` ou na raiz. Tire prints sempre
que precisar — são a evidência local do QA (a pasta é gitignored).

```bash
# agent-browser grava com caminho absoluto:
agent-browser screenshot /home/nilton/Projetos/nilton/go-snappymail/docs/prints/<contexto>/<nome>.png
# ex.: docs/prints/zimbra/…  ·  docs/prints/zimbra/qa8/…  ·  docs/prints/admin/…
```

Referências para comparar: webmail `https://192.168.56.30/` · admin `https://192.168.56.30:7071/`.

## Troubleshooting

- **Snapshot não existe** → `vagrant up` (instala) e depois `vagrant snapshot save zimbra-installed`.
- **VM não sobe / porta 53** → o provisionamento libera o `systemd-resolved`; ver `vagrant/zimbra/provision/01-dns.sh`.
- **TLS self-signed no webmail** → `insecure_skip_verify = true` nas seções `[imap]`/`[smtp]`.
- **Zimbra com serviço parado** → `vagrant ssh -c "sudo su - zimbra -c 'zmcontrol restart'"`.
- **Reset total do Zimbra** → `vagrant snapshot restore zimbra-installed` (descarta tudo desde o snapshot).
