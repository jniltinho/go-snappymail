<script setup lang="ts">
import { useAuthStore } from '../stores/auth'
import { useMailStore } from '../stores/mail'
import { useSettingsStore } from '../stores/settings'
import DropdownMenu from './DropdownMenu.vue'

import { computed } from 'vue'

const auth = useAuthStore()
const mail = useMailStore()
const settings = useSettingsStore()

const userLabel = computed(() => {
  const local = (auth.username || '').split('@')[0] || 'Account'
  return local.charAt(0).toUpperCase() + local.slice(1)
})

async function onSearch() {
  await mail.search()
}

function readRaw() {
  if (!mail.selectedUid) return
  window.open(
    `${API_BASE}/mail/${encodeURIComponent(mail.currentFolder)}/${mail.selectedUid}/raw`,
    '_blank',
  )
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

      <DropdownMenu :label="userLabel" btn-class="topbar-link font-bold" align-right>
        <button type="button" class="dd-item" @click="settings.toggleDark">
          {{ settings.darkMode ? 'Light theme' : 'Dark theme' }}
        </button>
        <button type="button" class="dd-item" @click="auth.logout">Sign out</button>
      </DropdownMenu>
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
        <DropdownMenu label="New message" btn-class="btn-new">
          <button type="button" class="dd-item" @click="mail.openCompose('new')">New message</button>
        </DropdownMenu>
      </div>

      <button type="button" class="tbtn" :disabled="!mail.selectedUid" @click="mail.openCompose('reply')">
        Reply
      </button>
      <button type="button" class="tbtn" :disabled="!mail.selectedUid" @click="mail.openCompose('replyall')">
        Reply to All
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

      <DropdownMenu label="Move" btn-class="tbtn">
        <button
          v-for="folder in mail.folders"
          :key="folder.name"
          type="button"
          class="dd-item"
          :disabled="!mail.selectedUid || folder.name === mail.currentFolder"
          @click="mail.selectedUid && mail.moveMessage(mail.selectedUid, folder.name)"
        >
          {{ folder.label === 'INBOX' ? 'Inbox' : folder.label }}
        </button>
      </DropdownMenu>

      <span class="toolbar-sep"></span>

      <DropdownMenu label="Actions" btn-class="tbtn">
        <button
          type="button"
          class="dd-item"
          :disabled="!mail.selectedUid"
          @click="mail.selectedUid && mail.setSeen(mail.selectedUid, !mail.selectedMessage?.seen)"
        >
          {{ mail.selectedMessage?.seen ? 'Mark unread' : 'Mark read' }}
        </button>
        <button
          type="button"
          class="dd-item"
          :disabled="!mail.selectedUid"
          @click="mail.selectedUid && mail.toggleFlag(mail.selectedUid, !mail.selectedMessage?.flagged)"
        >
          {{ mail.selectedMessage?.flagged ? 'Unflag' : 'Flag' }}
        </button>
        <button type="button" class="dd-item" :disabled="!mail.selectedUid" @click="mail.spamSelected">
          Mark as spam
        </button>
      </DropdownMenu>

      <div class="ml-auto flex items-center gap-2">
        <button type="button" class="tbtn" :disabled="!mail.selectedUid" @click="readRaw">
          Read More
        </button>
        <DropdownMenu label="View" btn-class="tbtn" align-right>
          <button type="button" class="dd-item" @click="settings.toggleDark">
            {{ settings.darkMode ? 'Light theme' : 'Dark theme' }}
          </button>
        </DropdownMenu>
      </div>
    </div>
  </header>
</template>
