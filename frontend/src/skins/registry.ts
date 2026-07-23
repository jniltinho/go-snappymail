/**
 * Re-exports from manifest.ts — prefer importing from ./manifest directly.
 * @see docs/skins.md
 */
export {
  SKIN_MANIFEST,
  SKIN_REGISTRY,
  DEFAULT_SKIN,
  normalizeSkinId,
  isSkinId,
  allSkinIds,
} from './manifest'
export type { SkinId, SkinMeta, SkinInfo, UIConfigResponse } from './manifest'
