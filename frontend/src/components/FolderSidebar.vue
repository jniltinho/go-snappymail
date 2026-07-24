<script setup lang="ts">
import { computed } from 'vue'
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

// Zimbra Classic folder order: Inbox, Sent, Drafts, Junk, Trash, then the rest.
const rank: Record<string, number> = { inbox: 0, sent: 1, drafts: 2, junk: 3, trash: 4 }

const ordered = computed(() =>
  [...mail.folders].sort((a, b) => {
    const ra = rank[a.iconType] ?? 5
    const rb = rank[b.iconType] ?? 5
    return ra !== rb ? ra - rb : a.label.localeCompare(b.label)
  }),
)

function prettyLabel(label: string): string {
  return label === 'INBOX' ? 'Inbox' : label
}
</script>

<template>
  <aside class="border-r border-line bg-panel overflow-y-auto min-h-0">
    <div class="side-header px-3 py-2">▼ Mail Folders</div>
    <button
      v-for="folder in ordered"
      :key="folder.name"
      type="button"
      class="side-item w-full text-left"
      :class="{ active: mail.currentFolder === folder.name }"
      :style="{ paddingLeft: `${12 + folder.depth * 14}px` }"
      @click="mail.selectFolder(folder.name)"
    >
      <span>{{ icons[folder.iconType] || icons.folder }}</span>
      <span class="truncate" :class="{ 'font-bold': folder.unseen }">
        {{ prettyLabel(folder.label) }}<template v-if="folder.unseen"> ({{ folder.unseen }})</template>
      </span>
    </button>
  </aside>
</template>
