/**
 * Skin registry — see docs/skins.md to add a new skin.
 * Run: make new-skin ID=<name>
 */
import type { SkinId, SkinMeta } from './types'

export const SKIN_REGISTRY: Record<SkinId, SkinMeta> = {
  snappymail: { id: 'snappymail', label: 'SnappyMail', ready: true },
  gmail: { id: 'gmail', label: 'Gmail', ready: false },
  outlook: { id: 'outlook', label: 'Outlook', ready: false },
}

export const DEFAULT_SKIN: SkinId = 'snappymail'

export function normalizeSkinId(raw: string | undefined): SkinId {
  const s = (raw || '').toLowerCase().trim()
  if (s === 'gmail' || s === 'google') return 'gmail'
  if (s === 'outlook' || s === 'office' || s === 'microsoft') return 'outlook'
  return 'snappymail'
}

export function isSkinId(id: string): id is SkinId {
  return id in SKIN_REGISTRY
}
