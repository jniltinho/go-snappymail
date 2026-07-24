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

      <form class="ml-auto flex items-center gap-1" @submit.prevent="onSearch">
        <input
          v-model="mail.searchQuery"
          type="search"
          placeholder="Search"
          class="topbar-search w-56 h-[24px] px-2 text-sm"
        />
        <button type="submit" class="topbar-link" title="Search">🔍</button>
      </form>

      <span class="text-xs hidden md:inline opacity-90 font-bold">{{ auth.username }} ▾</span>
      <button type="button" class="topbar-link" @click="settings.toggleDark">
        {{ settings.darkMode ? 'Light' : 'Dark' }}
      </button>
      <button type="button" class="topbar-link" @click="auth.logout">Logout</button>
    </div>

    <div class="tabstrip">
      <span class="tab active">Mail</span>
      <span class="tab disabled" title="Coming soon">Contacts</span>
      <span class="tab disabled" title="Coming soon">Calendar</span>
      <span class="tab disabled" title="Coming soon">Tasks</span>
      <span class="tab disabled" title="Coming soon">Preferences</span>
      <button type="button" class="tab-refresh ml-auto" title="Refresh" @click="mail.refresh">⟳</button>
    </div>

    <div class="actionbar flex items-center gap-2 px-3 py-2 border-b border-line">
      <div class="w-[204px] shrink-0">
        <button type="button" class="btn-new" @click="mail.openCompose('new')">New message ▾</button>
      </div>

      <button type="button" class="tbtn" :disabled="!mail.selectedUid" @click="mail.openCompose('reply')">
        Reply
      </button>
      <button type="button" class="tbtn" :disabled="!mail.selectedUid" @click="mail.openCompose('replyall')">
        Reply to all
      </button>
      <button type="button" class="tbtn" :disabled="!mail.selectedUid" @click="mail.openCompose('forward')">
        Forward
      </button>

      <span class="toolbar-sep"></span>

      <button type="button" class="tbtn" :disabled="!mail.selectedUid" @click="mail.archiveSelected">
        Archive
      </button>
      <button type="button" class="tbtn tbtn-danger" :disabled="!mail.selectedUid" @click="mail.deleteSelected">
        Delete
      </button>
      <button type="button" class="tbtn" :disabled="!mail.selectedUid" @click="mail.spamSelected">
        Spam
      </button>

      <span class="toolbar-sep"></span>

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
    </div>
  </header>
</template>
