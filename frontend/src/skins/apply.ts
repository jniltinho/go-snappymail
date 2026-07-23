import type { SkinId } from './types'

const SKIN_ATTR = 'data-skin'

/** Apply skin to the document root (drives CSS variable blocks per skin). */
export function applySkin(skinId: SkinId): void {
  document.documentElement.setAttribute(SKIN_ATTR, skinId)
}

export function currentSkin(): SkinId {
  const v = document.documentElement.getAttribute(SKIN_ATTR)
  if (v === 'gmail' || v === 'outlook' || v === 'snappymail') return v
  return 'snappymail'
}
