<script setup lang="ts">
import { computed } from 'vue'
import { useMailStore } from '../stores/mail'

const mail = useMailStore()

const fromLabel = computed(() => {
  const m = mail.selectedMessage
  if (!m) return ''
  if (m.from && m.fromEmail && m.from !== m.fromEmail) return `"${m.from}" <${m.fromEmail}>`
  return m.fromEmail || m.from
})

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
          <div class="msg-avatar" aria-hidden="true">
            <svg viewBox="0 0 48 48" width="48" height="48">
              <rect width="48" height="48" fill="#dbe9f5" />
              <circle cx="24" cy="17" r="8" fill="#5b87b5" />
              <path d="M8 44c1.5-10 8-15 16-15s14.5 5 16 15z" fill="#5b87b5" />
            </svg>
          </div>
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
              <span class="addr-chip truncate">{{ mail.selectedMessage.to }}</span>
            </div>
          </div>
        </div>
      </header>

      <div
        v-if="mail.selectedMessage.attachments?.length"
        class="attach-strip px-3 py-1.5 border-b border-line text-sm"
      >
        <span v-for="att in mail.selectedMessage.attachments" :key="att.part" class="mr-4">
          📎 {{ att.filename }} ({{ att.sizeLabel }})
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
