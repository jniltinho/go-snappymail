# Zimbra FOSS em VM Vagrant — instalação e configuração

Guia da VM de referência do projeto: **Zimbra FOSS 10.1.17.p2** (build comunitário
[maldua/zimbra-foss](https://github.com/maldua/zimbra-foss/releases)) em Ubuntu 24.04,
totalmente automatizada com Vagrant. Usada como referência visual viva da skin `zimbra`
do go-snappymail e como servidor de e-mail de laboratório.

---

## Visão geral

| Item | Valor |
|---|---|
| Build | `zcs-10.1.17_GA_4200002.UBUNTU24_64` ([downloads](https://maldua.github.io/zimbra-foss/downloads/)) |
| SO | Ubuntu 24.04 (box `bento/ubuntu-24.04`) |
| IP | `192.168.56.30` |
| FQDN | `mail.zimbra.test` |
| RAM/CPU | 10 GB / 4 vCPUs (mínimo recomendado pelo Zimbra: 8 GB) |
| Webmail | `https://192.168.56.30/` ou `https://mail.zimbra.test/` |
| Console admin | `https://mail.zimbra.test:9071/` |
| Admin | `admin@zimbra.test` / `Password1@` |
| Conta de teste | `nilton@linuxpro.com.br` / `Password1@` (domínio `linuxpro.com.br`) |

Estado atual do lab (após seed): domínios `linuxpro.com.br` (34 contas) e
`criarenet.com` (20 contas), todas com senha `Password1@`; caixa do
`nilton@linuxpro.com.br` com 230 mensagens de teste (pequenas, ~100KB,
anexos ~1MB e HTML); proxy em modo `both` (http **e** https) via
`zmprov ms mail.zimbra.test zimbraReverseProxyMailMode both`.

Arquivos: [`vagrant/zimbra/`](../vagrant/zimbra/) — `Vagrantfile`,
`provision/01-dns.sh`, `provision/02-zimbra.sh`.

---

## Subir a VM

```bash
cd vagrant/zimbra
vagrant up        # 1ª vez: baixa o tarball (236 MB) e instala — 20-40 min
```

No `/etc/hosts` do **host**:

```
192.168.56.30  mail.zimbra.test
```

O certificado é self-signed — aceite o aviso do browser. Para acessar o webmail
com o cliente **Classic** (o layout que a skin `zimbra` reproduz), use
`https://mail.zimbra.test/?client=advanced` (já é o default da VM).

Comandos úteis:

```bash
vagrant ssh                                   # shell na VM
vagrant ssh -c "sudo su - zimbra -c 'zmcontrol status'"   # status dos serviços
vagrant halt / vagrant up                     # parar / iniciar
vagrant destroy -f && vagrant up              # reinstalar do zero (usa o cache)
```

---

## Como a instalação automatizada funciona

### 1. DNS local (`provision/01-dns.sh`)

O Zimbra exige DNS funcional com registro **A** e **MX** para o domínio.
A VM resolve isso com `dnsmasq` local:

```
domain=zimbra.test
address=/mail.zimbra.test/192.168.56.30    # registro A
mx-host=zimbra.test,mail.zimbra.test,10    # registro MX
server=1.1.1.1                             # forward do resto
```

O stub do `systemd-resolved` é desabilitado (porta 53) e o `/etc/resolv.conf`
aponta para `127.0.0.1`.

### 2. Instalação unattended (`provision/02-zimbra.sh`)

Duas fases, como no instalador oficial:

1. **Pacotes** — `./install.sh -s --skip-activation-check` recebe as respostas
   por stdin (licença `y`, repositório `y`, pacotes default com `zimbra-dnscache`
   desabilitado — o dnsmasq já cumpre o papel).
2. **Configuração** — `/opt/zimbra/libexec/zmsetup.pl -c /tmp/zmsetup.cfg` com um
   arquivo de resposta completo (hostname, senhas LDAP, domínio `zimbra.test`,
   admin `admin@zimbra.test`, proxy/mailproxy, portas padrão 80/443/143/993).

O tarball fica em `vagrant/zimbra/cache/` (gitignored) e é reutilizado em
reprovisionamentos. Para trocar de versão, atualize as variáveis `ZCS`/`URL`
no topo do `02-zimbra.sh` com o link da [página de releases](https://github.com/maldua/zimbra-foss/releases).

---

## Administração: domínios e contas de e-mail

Todos os comandos rodam como usuário `zimbra` dentro da VM
(`vagrant ssh` → `sudo su - zimbra`), ou em uma linha:

```bash
vagrant ssh -c "sudo su - zimbra -c '<comando>'"
```

### Criar um domínio

```bash
zmprov cd linuxpro.com.br
```

### Criar uma conta

```bash
zmprov ca nilton@linuxpro.com.br 'Password1@' displayName 'Nilton'
```

### Outras operações comuns

```bash
zmprov -l gaa                                   # listar todas as contas
zmprov gad                                      # listar domínios
zmprov sp nilton@linuxpro.com.br 'NovaSenha1@'  # trocar senha
zmprov aaa nilton@linuxpro.com.br contato@linuxpro.com.br   # alias
zmprov ma nilton@linuxpro.com.br zimbraMailQuota 1073741824 # quota 1GB
zmprov da conta@dominio                         # apagar conta
zmprov dd dominio.com                           # apagar domínio (sem contas)
zmprov mc default zimbraPrefClientType advanced # cliente Classic como default
```

Tudo isso também existe na **console admin** (`https://mail.zimbra.test:9071/`):
*Manage → Domains / Accounts*.

### DNS para domínios extras no laboratório

A conta `nilton@linuxpro.com.br` recebe e envia e-mail **dentro** da VM sem DNS
extra (entrega local). Para que outros hosts do lab entreguem para
`linuxpro.com.br`, adicione o MX no dnsmasq da VM:

```bash
echo 'mx-host=linuxpro.com.br,mail.zimbra.test,10' | sudo tee -a /etc/dnsmasq.d/zimbra.conf
echo 'address=/mail.zimbra.test/192.168.56.30'     | sudo tee -a /etc/dnsmasq.d/zimbra.conf
sudo systemctl restart dnsmasq
```

Em produção seriam os registros públicos: `A mail.linuxpro.com.br`,
`MX linuxpro.com.br → mail.linuxpro.com.br`, além de SPF/DKIM/DMARC
(DKIM: `/opt/zimbra/libexec/zmdkimkeyutil -a -d linuxpro.com.br`).

---

## Acesso IMAP/SMTP (clientes externos, incl. go-snappymail)

| Protocolo | Porta | Nota |
|---|---|---|
| IMAPS | 993 | TLS (self-signed → `insecure_skip_verify` no lab) |
| IMAP | 143 | STARTTLS |
| SMTP | 25 / 587 | 587 com auth |

Exemplo de `config.toml` do go-snappymail apontando para a VM:

```toml
[imap]
host = "mail.zimbra.test"
port = 993
tls  = true
insecure_skip_verify = true   # lab: cert self-signed

[smtp]
host = "mail.zimbra.test"
port = 587
starttls = true
insecure_skip_verify = true
```

---

## Extrair assets do webclient (sem VM)

A skin harmony (usada como fonte da skin `zimbra`) pode ser minerada direto do
tarball, sem instalar nada:

```bash
tar xzf zcs-*.tgz --wildcards '*/packages/zimbra-mbox-webclient-war_*.deb'
dpkg-deb -x packages/zimbra-mbox-webclient-war_*.deb out/
ls out/opt/zimbra/jetty_base/webapps/zimbra/skins/harmony/
# manifest.xml  skin.css  skin.properties  img/
```

`skin.properties` contém as substituições de cor (AltC `#007CC3`, SelC
`lighten(AltC,60)` = `#99CAE7`, app row `#0087C3`, etc.) documentadas em
`openspec/changes/add-zimbra-skin/design.md`.

---

## Solução de problemas

- **`zmcontrol status` com serviço parado** — `sudo su - zimbra -c 'zmcontrol restart'` (leva ~2 min).
- **Instalador reclama de DNS** — confira `host -t MX zimbra.test 127.0.0.1` dentro da VM.
- **Reprovisionar sem re-download** — o cache em `vagrant/zimbra/cache/` é reaproveitado.
- **Memória** — com menos de 8 GB o zmsetup falha em `mailboxd`; o Vagrantfile aloca 10 GB.
- **Logs da instalação** — `/opt/zimbra/log/zmsetup.*.log` na VM.
