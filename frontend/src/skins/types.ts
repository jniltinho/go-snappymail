/** Known layout skin identifiers — extend when adding skins (docs/skins.md). */
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
