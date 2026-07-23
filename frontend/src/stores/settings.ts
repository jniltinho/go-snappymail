import { defineStore } from 'pinia'
import { ref, watch } from 'vue'

export const useSettingsStore = defineStore('settings', () => {
  const darkMode = ref(localStorage.getItem('gsn_dark') === '1')

  watch(darkMode, (v) => {
    localStorage.setItem('gsn_dark', v ? '1' : '0')
    document.documentElement.classList.toggle('dark', v)
  }, { immediate: true })

  function toggleDark(): void {
    darkMode.value = !darkMode.value
  }

  return { darkMode, toggleDark }
})
