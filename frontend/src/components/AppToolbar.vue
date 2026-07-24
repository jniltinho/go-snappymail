<script setup lang="ts">
import { useAuthStore } from '../stores/auth'
import { useMailStore } from '../stores/mail'
import { useSettingsStore } from '../stores/settings'
import DropdownMenu from './DropdownMenu.vue'

import { computed } from 'vue'

const auth = useAuthStore()
const mail = useMailStore()
const settings = useSettingsStore()

const tabs = [
  { id: 'mail', label: 'Mail' },
  { id: 'contacts', label: 'Contacts' },
  { id: 'calendar', label: 'Calendar' },
  { id: 'tasks', label: 'Tasks' },
  { id: 'preferences', label: 'Preferences' },
] as const

const userLabel = computed(() => {
  const local = (auth.username || '').split('@')[0] || 'Account'
  return local.charAt(0).toUpperCase() + local.slice(1)
})

async function onSearch() {
  await mail.search()
}

function showOriginal() {
  if (!mail.selectedUid) return
  window.open(
    `${API_BASE}/mail/${encodeURIComponent(mail.currentFolder)}/${mail.selectedUid}/raw`,
    '_blank',
  )
}

function printMessage() {
  window.print()
}

function editAsNew() {
  const m = mail.selectedMessage
  if (!m) return
  mail.openCompose('new')
  mail.cTo = m.to || ''
  mail.cSubject = m.subject
  mail.cBody = m.plainBody || ''
}
</script>

<template>
  <header class="shrink-0">
    <div class="topbar flex items-center gap-3 px-3">
      <div class="font-bold text-base tracking-tight whitespace-nowrap">go-snappymail</div>

      <form class="ml-auto relative" @submit.prevent="onSearch">
        <input
          v-model="mail.searchQuery"
          type="search"
          placeholder="Search"
          class="topbar-search w-56 h-[24px] pl-2 pr-7 text-sm"
        />
        <button type="submit" class="search-glass" title="Search">
          <svg viewBox="0 0 16 16" width="13" height="13" aria-hidden="true">
            <circle cx="6.5" cy="6.5" r="4.2" fill="none" stroke="currentColor" stroke-width="1.5" />
            <path d="M9.8 9.8L14 14" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" />
          </svg>
        </button>
      </form>

      <DropdownMenu :label="userLabel" btn-class="topbar-link font-bold" align-right>
        <button type="button" class="dd-item" @click="settings.toggleDark">
          {{ settings.darkMode ? 'Light theme' : 'Dark theme' }}
        </button>
        <button type="button" class="dd-item" @click="auth.logout">Sign out</button>
      </DropdownMenu>
    </div>

    <div class="tabstrip">
      <button
        v-for="t in tabs"
        :key="t.id"
        type="button"
        class="tab"
        :class="{ active: settings.activeTab === t.id }"
        @click="settings.activeTab = t.id"
      >
        {{ t.label }}
      </button>
      <button type="button" class="tab-refresh ml-auto" title="Refresh" @click="mail.refresh">⟳</button>
    </div>

    <div v-if="settings.activeTab === 'mail'" class="actionbar flex items-center gap-2 pr-3 py-1 border-b border-line">
      <div class="shrink-0" :style="{ width: `${settings.sideWidth}px`, paddingLeft: '10px' }">
        <DropdownMenu label="New Message" btn-class="btn-new" split @main="mail.openCompose('new')">
          <button type="button" class="dd-item" @click="mail.openCompose('new')">New Message</button>
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

      <DropdownMenu label="🏷" btn-class="tbtn tbtn-icon">
        <button type="button" class="dd-item" disabled>New Tag…</button>
        <button type="button" class="dd-item" disabled>No tags defined</button>
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
        <button type="button" class="dd-item" :disabled="!mail.selectedUid" @click="editAsNew">
          Edit as New
        </button>
        <button type="button" class="dd-item" :disabled="!mail.selectedUid" @click="showOriginal">
          Show Original
        </button>
        <button type="button" class="dd-item" :disabled="!mail.selectedUid" @click="printMessage">
          Print
        </button>
      </DropdownMenu>

      <div class="ml-auto flex items-center gap-2">
        <button type="button" class="tbtn" @click="mail.readNextUnread">
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
