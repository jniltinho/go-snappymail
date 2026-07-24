// Admin skin catalog. Mirrors the webmail skin manifest so both frontends share
// the same multi-skin model. The active skin comes from the server ([admin] skin)
// via the data-skin attribute on <html>; DEFAULT_SKIN is the fallback.
export interface AdminSkin {
  id: string
  label: string
  ready: boolean
}

export const ADMIN_SKINS: AdminSkin[] = [
  { id: 'serenity', label: 'Serenity', ready: true },
  { id: 'carbon', label: 'Carbon', ready: false },
]

export const DEFAULT_SKIN = 'serenity'

/** applySkin sets the active skin on <html>, falling back to the default. */
export function applySkin(id: string | null | undefined): void {
  const skin = ADMIN_SKINS.find((s) => s.id === id && s.ready)?.id ?? DEFAULT_SKIN
  document.documentElement.setAttribute('data-skin', skin)
}
