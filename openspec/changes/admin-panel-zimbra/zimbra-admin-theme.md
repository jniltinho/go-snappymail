# ZimbraAdmin theme tokens (referência para o `frontend-admin/`)

Valores extraídos do **skin `serenity`** do ZimbraAdmin real (VM `192.168.56.30:7071`,
`/opt/zimbra/jetty_base/webapps/zimbraAdmin/skins/serenity/skin.properties` +
`_base/base/skin.properties`). Assets crus (só referência, gitignored):
`docs/prints/zimbra-admin/assets/`. Usar estes valores no `tailwind.config` / `@theme`
do admin para o tema ficar **idêntico** ao legado. UI toda em inglês.

## Paleta (medida)

| Token (Zimbra) | Valor | Uso |
|---|---|---|
| `SerenityAltC` | `#0095D3` | **azul accent** (login gradient topo, seleção ativa, links) |
| `SelC` | `#C1DFFE` | seleção de linha/aba (azul claro) |
| `AppC` | `#CECECE` | chrome/cinza da app |
| `AltC` | `#DDD` | cinza secundário |
| `SerenityAppC` | `#EEEEEE` | base clara do serenity |
| `AppBg` | `#E7E9EE` | fundo da app |
| `viewC` | `#FCFCFC` | fundo de conteúdo/list view |
| `PanelColor` | `white` | painéis |
| `TxtC` | `#333` | texto |
| `AppPanelBorderC` | `#AAA` | bordas de painel (top/middle/footer) |
| `ButtonBorderColor` | `white` | borda de botão |
| Login bkgd | `lighten(#EEEEEE,55%)` → `darken(#EEEEEE,7%)` | gradiente claro da página de login |
| Login content box | `#0095D3` → `darken(#0095D3,33%)` | gradiente azul do topo do card de login |

## Tipografia (`_base/base/skin.properties`)

- **Família**: `"Segoe UI","Lucida Sans","Lucida Sans Unicode","Lucida Grande","DejaVu Sans",sans-serif`
- Tamanhos: normal **11px**, big 13px, bigger 15px bold, biggest 18px bold, small 10px
  - (o admin é mais compacto que o webmail — base 11px)

## Cantos / bordas / sombras

- Cantos: **3px** (abas `roundCorners(3px 3px 0 0)`, botões `roundCorners(3px)`, toast)
- Painéis: borda `1px solid #AAA` (AppTop/Middle/Footer)
- Login button: gradiente claro + `box-shadow 0 1px 3px rgba(50,50,50,.75)` + 3px + borda branca
- Toast: `box-shadow 0 0 10px <cor>` + cantos arredondados

## Layout (do console real)

- **Top bar**: faixa clara com banner "Administration" (marca texto no clone), busca central, `admin@… ▾`, refresh.
- **Árvore de navegação** (esquerda): Home / Monitor / Manage / Configure / Tools & Migration / Search.
- **Home**: "Overview" (Version, Servers, Accounts, Domains, COS) + "Runtime" + cards de setup — no clone só os counts reais (accounts/domains/aliases/admins), resto "n/a".
- **Content pane**: toolbar (New/Edit/Delete/…) + list view (colunas, seleção `#C1DFFE`, paginação).
- **Login**: card com topo em gradiente azul (`#0095D3`), labels à esquerda, botão Login canto inferior direito, página em gradiente cinza claro.

## Skins do admin (multi-skin, como o webmail)

O ZimbraAdmin real tem vários skins (`serenity`, `carbon`, `vami2`). O `frontend-admin/`
espelha isso — catálogo de skins com tokens por `[data-skin='<id>']`, `[admin] skin` na config:

| Skin | AppC | SelC (seleção) | Accent | Nota |
|---|---|---|---|---|
| **serenity** (default) | `#CECECE` | `#C1DFFE` | `#0095D3` | flat claro (o do console capturado) |
| **carbon** | `#cecece` | `#c4ddff` | (masthead escuro) | 2º skin; mesma estrutura, accent/masthead diferentes |
| vami2 / outras | — | — | — | espaço reservado |

Fonte base (`FontFamily-default`, 11px) e cantos 3px são comuns; cada skin varia paleta/accent.

## Divergências intencionais (clone)

- Marca textual (sem logo/banner Zimbra).
- Motor Vue 3 + TailwindCSS (não DWT/AJAX).
- **Inglês** (o console legado está em pt-BR pelo locale da VM).
- Cantos 3px mantidos (o admin serenity usa 3px — sem conflito com a regra do webmail, que é escopada por skin).

> Fluxo: capturar **prints** de cada tela (`docs/prints/zimbra-admin/`) + estes tokens → montar o tema Tailwind do admin → validar cada tela com `qa-frontend-cloner` contra `:7071`. **Sempre começar pelo backend** (API + testes) antes da UI.
