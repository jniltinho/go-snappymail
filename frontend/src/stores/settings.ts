import { defineStore } from 'pinia'
import { ref, watch } from 'vue'
import { applySkin } from '../skins/apply'
import { DEFAULT_SKIN, SKIN_REGISTRY, SKIN_MANIFEST, normalizeSkinId } from '../skins/registry'
import type { SkinId, UIConfigResponse, SkinInfo } from '../skins/manifest'

export const useSettingsStore = defineStore('settings', () => {
  const darkMode = ref(localStorage.getItem('gsn_dark') === '1')
  const activeTab = ref<'mail' | 'contacts' | 'calendar' | 'tasks' | 'preferences'>('mail')
  const sideWidth = ref(Number(localStorage.getItem('gsn_side_w')) || 190)
  const skin = ref<SkinId>(DEFAULT_SKIN)
  const availableSkins = ref<SkinInfo[]>([])
  const rowsPerPage = ref(50)
  const composeHTML = ref(true)

  function applyDarkClass(): void {
    document.documentElement.classList.toggle('dark', darkMode.value)
  }

  watch(darkMode, (v) => {
    localStorage.setItem('gsn_dark', v ? '1' : '0')
    applyDarkClass()
  }, { immediate: true })

  watch(skin, (id) => {
    applySkin(id)
  })

  function toggleDark(): void {
    darkMode.value = !darkMode.value
  }

  function initFromServer(config: UIConfigResponse): void {
    skin.value = normalizeSkinId(config.skin)
    if (config.skins?.length) {
      availableSkins.value = config.skins
    } else if (config.available_skins?.length) {
      availableSkins.value = config.available_skins.map((id) => ({
        id,
        label: SKIN_REGISTRY[normalizeSkinId(id)]?.label ?? id,
        ready: SKIN_REGISTRY[normalizeSkinId(id)]?.ready ?? false,
      }))
    } else {
      availableSkins.value = SKIN_MANIFEST.map(({ id, label, ready }) => ({ id, label, ready }))
    }
    rowsPerPage.value = config.rows_per_page || 50
    composeHTML.value = config.compose_html ?? true
    applySkin(skin.value)
    applyDarkClass()
  }

  function skinLabel(id: SkinId): string {
    return SKIN_REGISTRY[id]?.label ?? id
  }

  function skinReady(id: SkinId): boolean {
    return SKIN_REGISTRY[id]?.ready ?? false
  }

  return {
    darkMode,
    activeTab,
    sideWidth,
    skin,
    availableSkins,
    rowsPerPage,
    composeHTML,
    toggleDark,
    initFromServer,
    skinLabel,
    skinReady,
  }
})
