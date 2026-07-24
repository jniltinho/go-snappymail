# Zimbra Webapps — Folder / CSS / JS / HTML Structure

Structural manifest of the two front-end webapps on the reference install
(`/opt/zimbra/jetty_base/webapps/`). This documents **layout, file types, and
purpose** — it does not reproduce Zimbra's source. To read the actual code, use
the open-source repos in [08-source-repos.md](08-source-repos.md).

## Webapps deployed

| Webapp | Role |
| --- | --- |
| `zimbra` | The end-user **webmail** client (Zimbra Web Client / "Advanced" AjaxTk UI) |
| `zimbraAdmin` | The **admin console** (the UI we cloned) |
| `service` | The Java backend — SOAP/REST API, mail store, provisioning (no UI) |
| `zimlet` | Zimlet (plugin) hosting |

Both `zimbra` and `zimbraAdmin` are built on the **same JS UI framework** — the
**Zimbra Ajax Toolkit (AjaxTk)**, whose widget classes are prefixed `Dwt*`
(`DwtShell`, `DwtDialog`, `DwtTreeItem`, …). That's why the two share
look-and-feel and skins. When we inspected the admin login DOM we saw exactly
these: `DWT5.DwtShell` → `LoginScreen` → `.center` → `.contentBox`.

## `zimbra` — Webmail (sizes / file-type counts)

| Dir | Size | Contents |
| --- | --- | --- |
| `js/` | 41M | The AjaxTk framework + `ZmMail` app, split into bundles (see below) |
| `WEB-INF/` | 32M | JSPs, `.tag` files, web.xml, backend glue |
| `public/` | 22M | Static assets, login/error pages |
| `img/` | 4.4M | 763 png, 213 gif, 74 svg icons/sprites |
| `skins/` | 1.6M | **21 skins** (theme packs) — see below |
| `templates/` | 608K | 112 FreeMarker `.ftl` + 53 `.tag` UI templates |
| `css/` | 504K | 75 stylesheets |
| `h/` | 132K | The "standard" (HTML, non-AJAX) client |
| `portals/`, `downloads/`, `qunit/` | — | Portals, installers, tests |

File types (webmail): `876 js, 763 png, 500 properties (i18n), 213 gif,
168 bcmap (pdf.js fonts), 112 ftl, 75 css, 74 svg, 62 jsp, 53 tag,
42 zgz, 42 appcache, 27 xml, 17 html`.

### JS bundles (webmail `js/`)

The client is compiled into feature bundles, each shipped four ways:
`X.js` (dev), `X_all.js` (concatenated), `X_all.js.zgz` (pre-gzipped),
`X.appcache` (offline manifest). Bundles:

`Boot, Startup1_1/1_2, Startup2, Ajax, MailCore, Mail, CalendarCore, Calendar,
CalendarAppt, ContactsCore, Contacts, Preferences(Core), Briefcase(Core),
Tasks(Core), Docs(Preview), Share, Portal, Zimlet(App), TinyMCE, JQuery,
TwoFactor, QRCode, PasswordRecovery, Voicemail, ImportExport, Crypt, …`

Source subtrees: `js/ajax` (the DWT toolkit), `js/zimbra` (Zimbra common),
`js/zimbraMail` (the mail app).

### CSS (webmail `css/`)

Core: `zm.css`, `dwt.css` (widget toolkit), `common.css`, `login.css` +
`zlogin.css` (login screen), `msgview.css`, `editor.css`, `spellcheck.css`.
Device variants: `iphone*.css`, `ipad.css`, `wm6*.css`, `zmobile*.css`,
`zhtml.css`, `xlite.css`. (Desktop theming proper lives in `skins/`, not here.)

## `zimbraAdmin` — Admin console (sizes / file-type counts)

| Dir | Size | Contents |
| --- | --- | --- |
| `WEB-INF/` | 29M | JSPs, tags, backend glue |
| `help/` | 19M | 575 `.htm` help pages + 12 PDFs |
| `js/` | 17M | AjaxTk + the `Admin` app bundles |
| `skins/` | 848K | **4 skins**: `_base, carbon, serenity, vami2` |
| `img/` | 820K | 130 png, 105 gif |
| `css/` | 500K | 41 stylesheets |
| `templates/` | 80K | UI templates |

File types (admin): `575 htm, 563 js, 315 properties, 130 png, 105 gif,
41 css, 16 jsp, 15 jpg, 14 xml, 11 jar, 11 html`.

JS bundles (admin `js/`): `Boot, Ajax, Admin, Zimbra, JQuery, XForms, Chartjs,
Clipboard, Debug` (same four-file shipping pattern). Subtrees: `js/ajax`,
`js/zimbra`, `js/zimbraAdmin`.

## Skins — the theme system (what we cloned)

A skin is a self-contained theme pack. The `serenity` admin skin (the one we
reproduced) contains:

```
skins/serenity/
  skin.properties   # 225 lines of theme tokens (colors, fonts, corners)
  skin.css          # generated CSS from the tokens
  skin.html         # skin HTML frame
  skin.js           # skin behavior
  manifest.xml      # skin manifest
  img/
    Decoration.png
    gradient-bg.png
    login-bg.png
    images.css       # sprite map
  logos/
    LoginBanner.png  # <- the login logo we extracted
    AppBanner.png
    AppBannerWhite.png
    AboutBanner.png
```

`skin.properties` is the source of the theme tokens we mirrored into
`frontend-admin/src/skins/serenity/theme.css` (accent `#0095D3`, selection
`#C1DFFE`, chrome `#CECECE`, 3px corners, Segoe UI 11px). See the theme mapping
in `openspec/changes/admin-panel-zimbra/zimbra-admin-theme.md`.

> The 21 webmail skins vs 4 admin skins explains why our webmail clone has more
> skin choices than the admin clone. Our multi-skin architecture mirrors this
> (`[data-skin]` per skin), currently `serenity` + a `carbon` stub.

## Server-side templating

- **FreeMarker** (`.ftl`, 112 in webmail) — server-rendered fragments.
- **JSP + tag libraries** (`.jsp`, `.tag`, `zm-taglib.jar`) — page assembly,
  login pages, bootstrapping.
- **`.properties`** (500 webmail / 315 admin) — i18n message catalogs, one set
  per locale.
