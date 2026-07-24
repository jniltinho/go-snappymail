<script setup lang="ts">
import { computed } from 'vue'
import { useMailStore } from '../stores/mail'

const mail = useMailStore()

// Classic gray/graphite sprites (Zimbra palette: gray-tan folders, graphite glyphs).
type IconPart = { d: string; fill?: string; stroke?: string; sw?: string }
const FOLDER_BODY: IconPart[] = [
  { d: 'M1 5.2V3.6h4.4l1.1 1.6H15v9H1z', fill: '#b3ac9c' },
  { d: 'M1 6.8h14', stroke: '#9c9484', sw: '0.8' },
]
const icons: Record<string, IconPart[]> = {
  inbox: [
    { d: 'M2 8.5l1.6-4h8.8L14 8.5v5H2z', fill: '#9a9a9a' },
    { d: 'M2 8.5h3.4l1 1.6h3.2l1-1.6H14', stroke: '#f2f2f2', sw: '0.9' },
    { d: 'M8 2.6v3.6M6.6 4.8L8 6.4l1.4-1.6', stroke: '#4a4a4a', sw: '1.3' },
  ],
  sent: [
    ...FOLDER_BODY,
    { d: 'M5.6 8.4h4.8v3.2H5.6z', fill: '#f5f5f0' },
    { d: 'M5.6 8.4l2.4 1.7 2.4-1.7', stroke: '#8a8578', sw: '0.8' },
  ],
  drafts: [
    ...FOLDER_BODY,
    { d: 'M6 12.2l.5-1.9 2.9-2.9 1.4 1.4-2.9 2.9z', fill: '#e8c94a' },
    { d: 'M6 12.2l.5-1.9 1.4 1.4z', fill: '#d0452e' },
  ],
  junk: [
    ...FOLDER_BODY,
    { d: 'M8.5 7.2a2.5 2.5 0 100 5 2.5 2.5 0 000-5z', stroke: '#c23b2e', sw: '1.3' },
    { d: 'M6.8 11.4l3.4-3.4', stroke: '#c23b2e', sw: '1.3' },
  ],
  trash: [
    { d: 'M4.2 5.5h7.6l-.8 8H5z', fill: '#9a9a9a' },
    { d: 'M3.4 5.5h9.2M6.4 5.5V3.9h3.2v1.6', stroke: '#7d7d7d', sw: '1.4' },
    { d: 'M6.6 7.4v4.2M8 7.4v4.2M9.4 7.4v4.2', stroke: '#f2f2f2', sw: '0.9' },
  ],
  folder: FOLDER_BODY,
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
          v-for="(part, i) in icons[folder.iconType] || icons.folder"
          :key="i"
          :d="part.d"
          :fill="part.fill || 'none'"
          :stroke="part.stroke"
          :stroke-width="part.sw"
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
