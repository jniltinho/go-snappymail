import { defineStore } from 'pinia'
import { ref, watch } from 'vue'
import { applySkin } from '../skins/apply'
import { DEFAULT_SKIN, SKIN_REGISTRY, normalizeSkinId } from '../skins/registry'
import type { SkinId, UIConfigResponse } from '../skins/types'

export const useSettingsStore = defineStore('settings', () => {
  const darkMode = ref(localStorage.getItem('gsn_dark') === '1')
  const skin = ref<SkinId>(DEFAULT_SKIN)
  const availableSkins = ref<string[]>(['snappymail', 'gmail', 'outlook'])
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
    availableSkins.value = config.available_skins?.length
      ? config.available_skins
      : Object.keys(SKIN_REGISTRY)
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
