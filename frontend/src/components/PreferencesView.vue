<script setup lang="ts">
import { useSettingsStore } from '../stores/settings'
import type { SkinId } from '../skins/manifest'
import MiniCalendar from './MiniCalendar.vue'

const settings = useSettingsStore()

const sections = [
  'General',
  'Accounts',
  'Mail',
  'Filters',
  'Signatures',
  'Out of Office',
  'Trusted Addresses',
  'Contacts',
  'Calendar',
  'Sharing',
  'Notifications',
  'Import / Export',
  'Shortcuts',
  'Zimlets',
]

function onThemeChange(e: Event) {
  settings.skin = (e.target as HTMLSelectElement).value as SkinId
}

function done() {
  settings.activeTab = 'mail'
}
</script>

<template>
  <div class="flex flex-col flex-1 min-h-0">
    <div class="actionbar flex items-center gap-2 px-3 py-1 border-b border-line">
      <button type="button" class="tbtn" @click="done">Save</button>
      <button type="button" class="tbtn" @click="done">Cancel</button>
      <span class="toolbar-sep"></span>
      <button type="button" class="tbtn" disabled>Undo Changes</button>
    </div>

    <div class="grid flex-1 min-h-0" style="grid-template-columns: 190px 1fr">
      <aside class="border-r border-line bg-panel overflow-y-auto flex flex-col">
        <div class="side-header px-3 py-2">▼ Preferences</div>
        <button
          v-for="(s, i) in sections"
          :key="s"
          type="button"
          class="side-item w-full text-left"
          :class="{ active: i === 0 }"
        >
          <span>⚙</span><span>{{ s }}</span>
        </button>
        <MiniCalendar />
      </aside>

      <section class="bg-panel min-h-0 overflow-y-auto px-4 py-3">
        <div class="pref-bar">Sign in</div>
        <div class="pref-row">
          <span class="pref-label">Password:</span>
          <button type="button" class="tbtn" disabled>Change Password</button>
        </div>

        <div class="pref-bar">Appearance</div>
        <div class="pref-row">
          <span class="pref-label">Theme:</span>
          <select class="pref-select" :value="settings.skin" @change="onThemeChange">
            <option v-for="s in settings.availableSkins" :key="s.id" :value="s.id">
              {{ s.id === 'zimbra' ? 'Harmony' : s.label }}
            </option>
          </select>
        </div>
        <div class="pref-row">
          <span class="pref-label">Dark mode:</span>
          <label class="text-sm flex items-center gap-1.5">
            <input type="checkbox" :checked="settings.darkMode" @change="settings.toggleDark" />
            Enable dark theme
          </label>
        </div>
        <div class="pref-row">
          <span class="pref-label">Font:</span>
          <select class="pref-select" disabled><option>Standard</option></select>
        </div>
        <div class="pref-row">
          <span class="pref-label">Display Font Size:</span>
          <select class="pref-select" disabled><option>Normal</option></select>
        </div>

        <div class="pref-bar">Time Zone and Language</div>
        <div class="pref-row">
          <span class="pref-label">Time Zone:</span>
          <select class="pref-select" disabled><option>GMT -03:00 Brasilia</option></select>
        </div>
        <div class="pref-row">
          <span class="pref-label">Language:</span>
          <select class="pref-select" disabled><option>English (United States)</option></select>
        </div>

        <div class="pref-bar">Search</div>
        <div class="pref-row">
          <span class="pref-label">Search Folders:</span>
          <label class="text-sm flex items-center gap-1.5">
            <input type="checkbox" disabled /> Include Spam Folder in Searches
          </label>
        </div>
        <div class="pref-row">
          <span class="pref-label"></span>
          <label class="text-sm flex items-center gap-1.5">
            <input type="checkbox" disabled /> Include Trash Folder in Searches
          </label>
        </div>

        <div class="pref-bar">Other Settings</div>
        <div class="pref-row">
          <span class="pref-label">Selection:</span>
          <label class="text-sm flex items-center gap-1.5">
            <input type="checkbox" disabled /> Display checkboxes to select items in lists
          </label>
        </div>
      </section>
    </div>
  </div>
</template>
