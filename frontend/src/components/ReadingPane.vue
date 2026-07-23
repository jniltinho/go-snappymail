<script setup lang="ts">
import { computed } from 'vue'
import { useMailStore } from '../stores/mail'

const mail = useMailStore()

const initial = computed(() => {
  const m = mail.selectedMessage
  const src = m?.from || m?.fromEmail || '?'
  return src.trim().charAt(0).toUpperCase() || '?'
})
</script>

<template>
  <section class="bg-panel overflow-y-auto min-h-0 flex flex-col">
    <template v-if="mail.selectedMessage">
      <header class="px-4 py-3 border-b border-line">
        <div class="flex gap-3">
          <div class="msg-avatar">{{ initial }}</div>
          <div class="flex-1 min-w-0">
            <div class="flex items-start justify-between gap-3">
              <h1 class="text-base font-semibold truncate">{{ mail.selectedMessage.subject }}</h1>
              <span class="text-xs text-ink-mute whitespace-nowrap">{{ mail.selectedMessage.date }}</span>
            </div>
            <div class="mt-1 flex items-center gap-2">
              <span class="hdr-label">From:</span>
              <span class="addr-chip">{{ mail.selectedMessage.from || mail.selectedMessage.fromEmail }}</span>
            </div>
            <div v-if="mail.selectedMessage.to" class="mt-1 flex items-center gap-2">
              <span class="hdr-label">To:</span>
              <span class="addr-chip">{{ mail.selectedMessage.to }}</span>
            </div>
          </div>
        </div>
      </header>

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

      <ul v-if="mail.selectedMessage.attachments?.length" class="px-4 pb-4 text-sm">
        <li v-for="att in mail.selectedMessage.attachments" :key="att.part" class="py-1">
          📎 {{ att.filename }} ({{ att.sizeLabel }})
        </li>
      </ul>
    </template>
    <div v-else class="flex-1 grid place-items-center text-ink-mute text-sm">
      To view a message, click on it.
    </div>
  </section>
</template>
