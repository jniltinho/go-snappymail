# VM Zimbra FOSS — referência visual da skin `zimbra`

VM Ubuntu 24.04 com **Zimbra FOSS 10.1.17.p2** ([maldua build](https://maldua.github.io/zimbra-foss/downloads/)),
usada como referência viva para a skin `zimbra` do go-snappymail (client Classic/harmony).

## Uso

```bash
cd vagrant/zimbra
vagrant up          # instalação unattended, 20-40 min (baixa 236MB na 1ª vez)
```

No `/etc/hosts` do host:

```
192.168.56.30  mail.zimbra.test
```

| Serviço | URL | Credenciais |
|---|---|---|
| Webmail | https://mail.zimbra.test/ | `admin@zimbra.test` / `Password1@` |
| Admin console | https://mail.zimbra.test:9071/ | idem |

## Notas

- `cache/` guarda o tarball (236MB, gitignored) para reprovisionar sem re-download.
- DNS interno via dnsmasq (A + MX de `zimbra.test` → `192.168.56.30`); a opção
  `zimbra-dnscache` do instalador fica desabilitada.
- Client Classic (advanced) é o default (`zimbraPrefClientType advanced`) — é o
  layout que a skin `zimbra` reproduz.
- Assets do webclient (skin harmony, `skin.properties`) podem ser extraídos do
  tarball sem VM: pacote `zimbra-mbox-webclient-war_*.deb` →
  `opt/zimbra/jetty_base/webapps/zimbra/skins/harmony/`.
