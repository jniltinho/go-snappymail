<script setup lang="ts">
import { computed } from 'vue'
import { useMailStore } from '../stores/mail'
import personAvatar from '../assets/person48.png'

const mail = useMailStore()

const fromLabel = computed(() => {
  const m = mail.selectedMessage
  if (!m) return ''
  if (m.from && m.fromEmail && m.from !== m.fromEmail) return `"${m.from}" <${m.fromEmail}>`
  return m.fromEmail || m.from
})

// `Name <email>` → `"Name" <email>` (Zimbra chip format)
const toLabel = computed(() =>
  (mail.selectedMessage?.to || '').replace(/(^|, )([^"<,][^<,]*?) </g, '$1"$2" <'),
)

function fmtSize(label: string): string {
  return label.replace(/(\d)([KMG]?B)$/, '$1 $2')
}

function attachmentURL(part: number): string {
  return `${API_BASE}/mail/${encodeURIComponent(mail.currentFolder)}/${mail.selectedUid}/attachment/${part}`
}
</script>

<template>
  <section class="bg-panel overflow-y-auto min-h-0 flex flex-col">
    <template v-if="mail.selectedMessage">
      <div class="subject-bar flex items-center justify-between px-3">
        <h1 class="subject-title truncate">{{ mail.selectedMessage.subject }}</h1>
        <span class="text-xs text-ink-sub whitespace-nowrap">1 message</span>
      </div>

      <header class="px-4 py-3 border-b border-line">
        <div class="flex gap-3">
          <img class="msg-avatar" :src="personAvatar" width="48" height="48" alt="" aria-hidden="true" />
          <div class="flex-1 min-w-0">
            <div class="flex items-start justify-between gap-3">
              <div class="mt-1 flex items-center gap-2 min-w-0">
                <span class="hdr-label">From:</span>
                <span class="addr-chip truncate">{{ fromLabel }}</span>
              </div>
              <span class="text-xs text-ink-sub whitespace-nowrap mt-1">{{
                mail.selectedMessage.dateFull
              }}</span>
            </div>
            <div v-if="mail.selectedMessage.to" class="mt-1 flex items-center gap-2">
              <span class="hdr-label">To:</span>
              <span class="addr-chip truncate">{{ toLabel }}</span>
            </div>
          </div>
        </div>
      </header>

      <div
        v-if="mail.selectedMessage.attachments?.length"
        class="attach-strip px-3 py-1.5 border-b border-line text-sm"
      >
        <span v-for="att in mail.selectedMessage.attachments" :key="att.part" class="mr-4">
          📎 {{ att.filename }} ({{ fmtSize(att.sizeLabel) }})
          <a class="attach-link ml-1" :href="attachmentURL(att.part)" target="_blank">Download</a>
        </span>
      </div>

      <div class="flex-1 p-4 prose prose-sm max-w-none dark:prose-invert">
        <div
          v-if="mail.selectedMessage.htmlBody"
          v-html="mail.selectedMessage.htmlBody"
        />
        <pre v-else-if="mail.selectedMessage.plainBody" class="whitespace-pre-wrap text-sm">{{
          mail.selectedMessage.plainBody
        }}</pre>
        <p v-else class="text-ink-mute text-sm">Loading message…</p>
      </div>
    </template>
    <div v-else class="flex-1 grid place-items-center text-ink-mute text-sm">
      To view a message, click on it.
    </div>
  </section>
</template>
