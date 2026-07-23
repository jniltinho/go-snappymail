/** Known layout skin identifiers (must match server internal/ui/skins.go). */
export type SkinId = 'snappymail' | 'gmail' | 'outlook'

export interface SkinMeta {
  id: SkinId
  label: string
  /** When false, skin CSS exists but full layout is not implemented yet. */
  ready: boolean
}

export interface UIConfigResponse {
  skin: string
  available_skins: string[]
  rows_per_page: number
  datetime_format: string
  compose_html: boolean
}
