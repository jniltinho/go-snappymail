<script setup lang="ts">
import { useAuthStore } from '../stores/auth'
import { useMailStore } from '../stores/mail'
import { useSettingsStore } from '../stores/settings'

const auth = useAuthStore()
const mail = useMailStore()
const settings = useSettingsStore()

async function onSearch() {
  await mail.search()
}
</script>

<template>
  <header class="flex items-center gap-2 px-3 py-2 border-b border-line bg-panel-2 shrink-0">
    <button type="button" class="tbtn" @click="mail.refresh">Refresh</button>
    <button
      type="button"
      class="tbtn"
      :disabled="!mail.selectedUid"
      @click="mail.selectedUid && mail.toggleFlag(mail.selectedUid, !mail.selectedMessage?.flagged)"
    >
      Flag
    </button>
    <button type="button" class="tbtn tbtn-danger" :disabled="!mail.selectedUid" @click="mail.deleteSelected">
      Delete
    </button>

    <form class="flex-1 flex gap-2 max-w-md ml-4" @submit.prevent="onSearch">
      <input
        v-model="mail.searchQuery"
        type="search"
        placeholder="Search this folder…"
        class="flex-1 h-[26px] px-2 border border-line bg-panel text-sm"
      />
      <button type="submit" class="tbtn">Search</button>
    </form>

    <button type="button" class="tbtn ml-auto" @click="settings.toggleDark">
      {{ settings.darkMode ? 'Light' : 'Dark' }}
    </button>
    <span class="text-xs text-ink-mute hidden md:inline">{{ auth.username }}</span>
    <button type="button" class="tbtn" @click="auth.logout">Logout</button>
  </header>
</template>
