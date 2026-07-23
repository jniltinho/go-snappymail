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
  <header class="shrink-0">
    <div class="topbar flex items-center gap-3 px-3">
      <div class="font-bold text-base tracking-tight whitespace-nowrap">go-snappymail</div>

      <form class="flex-1 flex justify-center gap-2 max-w-xl mx-auto" @submit.prevent="onSearch">
        <input
          v-model="mail.searchQuery"
          type="search"
          placeholder="Search"
          class="topbar-search flex-1 h-[26px] px-2 text-sm"
        />
        <button type="submit" class="tbtn">Search</button>
      </form>

      <button type="button" class="topbar-link" @click="settings.toggleDark">
        {{ settings.darkMode ? 'Light' : 'Dark' }}
      </button>
      <span class="text-xs hidden md:inline opacity-90">{{ auth.username }}</span>
      <button type="button" class="topbar-link" @click="auth.logout">Logout</button>
    </div>

    <div class="tabstrip">
      <span class="tab">Mail</span>
    </div>

    <div class="actionbar flex items-center gap-2 px-3 py-2 border-b border-line">
      <button type="button" class="tbtn" @click="mail.refresh">Refresh</button>
      <button
        type="button"
        class="tbtn"
        :disabled="!mail.selectedUid"
        @click="mail.selectedUid && mail.toggleFlag(mail.selectedUid, !mail.selectedMessage?.flagged)"
      >
        Flag
      </button>
      <button
        type="button"
        class="tbtn"
        :disabled="!mail.selectedUid"
        @click="mail.selectedUid && mail.setSeen(mail.selectedUid, !mail.selectedMessage?.seen)"
      >
        {{ mail.selectedMessage?.seen ? 'Mark unread' : 'Mark read' }}
      </button>
      <button type="button" class="tbtn" :disabled="!mail.selectedUid" @click="mail.archiveSelected">
        Archive
      </button>
      <button type="button" class="tbtn tbtn-danger" :disabled="!mail.selectedUid" @click="mail.deleteSelected">
        Delete
      </button>
    </div>
  </header>
</template>
