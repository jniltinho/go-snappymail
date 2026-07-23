/**
 * Skin manifest — keep in sync with internal/ui/skins.go (catalog).
 * Run: make validate-skins after edits.
 *
 * manifest-begin
 */
export const SKIN_MANIFEST = [
  {
    id: 'snappymail',
    label: 'SnappyMail',
    ready: true,
    aliases: ['default', 'snappymail-default', 'theme-default'],
  },
  {
    id: 'gmail',
    label: 'Gmail',
    ready: false,
    aliases: ['google'],
  },
  {
    id: 'outlook',
    label: 'Outlook',
    ready: true,
    aliases: ['office', 'microsoft'],
  },
  {
    id: 'carbonio',
    label: 'Carbonio',
    ready: true,
    aliases: ['zextras'],
  },
] as const // manifest-end

export type SkinId = (typeof SKIN_MANIFEST)[number]['id']

export interface SkinMeta {
  id: SkinId
  label: string
  ready: boolean
  aliases: readonly string[]
}

export interface SkinInfo {
  id: string
  label: string
  ready: boolean
}

export interface UIConfigResponse {
  skin: string
  skins?: SkinInfo[]
  available_skins: string[]
  rows_per_page: number
  datetime_format: string
  compose_html: boolean
}

export const DEFAULT_SKIN: SkinId = 'snappymail'

function buildRegistry(): Record<SkinId, SkinMeta> {
  const out = {} as Record<SkinId, SkinMeta>
  for (const s of SKIN_MANIFEST) {
    out[s.id] = { id: s.id, label: s.label, ready: s.ready, aliases: s.aliases }
  }
  return out
}

/** Lookup table built from manifest. */
export const SKIN_REGISTRY = buildRegistry()

const aliasMap = new Map<string, SkinId>(
  SKIN_MANIFEST.flatMap((s) =>
    [s.id, ...s.aliases].map((a) => [a.toLowerCase(), s.id] as const),
  ),
)

export function normalizeSkinId(raw: string | undefined): SkinId {
  const key = (raw || '').toLowerCase().trim()
  return aliasMap.get(key) ?? DEFAULT_SKIN
}

export function isSkinId(id: string): id is SkinId {
  return id in SKIN_REGISTRY
}

export function allSkinIds(): SkinId[] {
  return SKIN_MANIFEST.map((s) => s.id)
}
