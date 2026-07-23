import { axios, primeCsrf } from '../api/client'
import { applySkin } from './apply'
import { DEFAULT_SKIN, normalizeSkinId } from './registry'
import type { SkinId, UIConfigResponse } from './types'

export interface BootstrapUIResult {
  skin: SkinId
  config: UIConfigResponse
}

/** Load server UI defaults before the app mounts (skin from config.toml). */
export async function bootstrapUI(): Promise<BootstrapUIResult> {
  await primeCsrf()
  try {
    const res = await axios.get<UIConfigResponse>(`${API_BASE}/ui/config`)
    const skin = normalizeSkinId(res.data.skin)
    applySkin(skin)
    return { skin, config: res.data }
  } catch {
    applySkin(DEFAULT_SKIN)
    return {
      skin: DEFAULT_SKIN,
      config: {
        skin: DEFAULT_SKIN,
        available_skins: ['snappymail', 'gmail', 'outlook'],
        rows_per_page: 50,
        datetime_format: '02/01/2006 15:04',
        compose_html: true,
      },
    }
  }
}
