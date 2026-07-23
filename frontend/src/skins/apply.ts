import { normalizeSkinId, isSkinId, DEFAULT_SKIN } from './manifest'
import type { SkinId } from './manifest'

const SKIN_ATTR = 'data-skin'

/** Apply skin to the document root (drives CSS variable blocks per skin). */
export function applySkin(skinId: SkinId): void {
  document.documentElement.setAttribute(SKIN_ATTR, skinId)
}

export function currentSkin(): SkinId {
  const v = document.documentElement.getAttribute(SKIN_ATTR)
  if (v && isSkinId(v)) return v
  return DEFAULT_SKIN
}

/** Re-export for callers that resolve raw server strings. */
export { normalizeSkinId }
