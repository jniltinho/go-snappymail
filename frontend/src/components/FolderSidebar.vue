<script setup lang="ts">
import { computed } from 'vue'
import { useMailStore } from '../stores/mail'

const mail = useMailStore()

// Filled folder sprites (Zimbra-style): blue folder body + white overlay glyph.
type IconPart = { d: string; fill?: string; stroke?: string; sw?: string }
const FOLDER_BODY: IconPart[] = [
  { d: 'M1 5.2V3.6h4.4l1.1 1.6H15v9H1z', fill: '#7ba0c4' },
]
const icons: Record<string, IconPart[]> = {
  inbox: [
    ...FOLDER_BODY,
    { d: 'M8 6.8v3.2M6.5 8.6L8 10.2l1.5-1.6', stroke: '#ffffff', sw: '1.3' },
  ],
  sent: [
    ...FOLDER_BODY,
    { d: 'M6 9.6h3.6M8.4 7.9l1.7 1.7-1.7 1.7', stroke: '#ffffff', sw: '1.3' },
  ],
  drafts: [
    ...FOLDER_BODY,
    { d: 'M6 12l.5-2 3-3 1.5 1.5-3 3z', fill: '#ffffff' },
  ],
  junk: [
    ...FOLDER_BODY,
    { d: 'M8.5 7.4a2.3 2.3 0 100 4.6 2.3 2.3 0 000-4.6z', stroke: '#ffffff', sw: '1.2' },
    { d: 'M6.9 11.3l3.2-3.2', stroke: '#ffffff', sw: '1.2' },
  ],
  trash: [
    { d: 'M4 5.5h8l-.9 8.2H4.9z', fill: '#7ba0c4' },
    { d: 'M3.2 5.5h9.6M6.3 5.5V3.8h3.4v1.7', stroke: '#7ba0c4', sw: '1.4' },
    { d: 'M6.5 7.5v4M8 7.5v4M9.5 7.5v4', stroke: '#ffffff', sw: '0.9' },
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
