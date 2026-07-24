<script setup lang="ts">
import { computed } from 'vue'
import { useMailStore } from '../stores/mail'

const mail = useMailStore()

// Monochrome line-art folder icons (Zimbra-style sprites), stroke inherits currentColor.
const iconPaths: Record<string, string> = {
  inbox: 'M2 9l3-6h6l3 6v4H2V9zm0 0h4l1 2h2l1-2h4',
  sent: 'M2 8l12-5-4 11-3-4-5-2z',
  drafts: 'M3 13h10M4 11l7-7 2 2-7 7H4v-2z',
  junk: 'M8 2l6 11H2L8 2zm0 4v3m0 2v1',
  trash: 'M3 5h10M6 5V3h4v2M4 5l1 8h6l1-8',
  folder: 'M2 4h5l1 2h6v7H2V4z',
}

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
      <svg class="side-icon" viewBox="0 0 16 16" width="15" height="15" aria-hidden="true">
        <path
          :d="iconPaths[folder.iconType] || iconPaths.folder"
          fill="none"
          stroke="currentColor"
          stroke-width="1.3"
          stroke-linejoin="round"
          stroke-linecap="round"
        />
      </svg>
      <span class="truncate" :class="{ 'font-bold': folder.unseen }">
        {{ prettyLabel(folder.label) }}<template v-if="folder.unseen"> ({{ folder.unseen }})</template>
      </span>
    </button>
  </aside>
</template>
