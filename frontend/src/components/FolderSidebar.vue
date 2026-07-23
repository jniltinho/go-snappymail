<script setup lang="ts">
import { useMailStore } from '../stores/mail'

const mail = useMailStore()

const icons: Record<string, string> = {
  inbox: '📥',
  sent: '📤',
  drafts: '📝',
  trash: '🗑',
  junk: '⚠',
  folder: '📁',
}
</script>

<template>
  <aside class="border-r border-line bg-panel overflow-y-auto min-h-0">
    <div class="px-3 py-2 text-xs font-semibold uppercase tracking-wide text-ink-mute">
      Folders
    </div>
    <button
      v-for="folder in mail.folders"
      :key="folder.name"
      type="button"
      class="side-item w-full text-left"
      :class="{ active: mail.currentFolder === folder.name }"
      :style="{ paddingLeft: `${12 + folder.depth * 14}px` }"
      @click="mail.selectFolder(folder.name)"
    >
      <span>{{ icons[folder.iconType] || icons.folder }}</span>
      <span class="truncate flex-1">{{ folder.label }}</span>
      <span v-if="folder.unseen" class="text-xs font-bold text-accent-2">{{ folder.unseen }}</span>
    </button>
  </aside>
</template>
