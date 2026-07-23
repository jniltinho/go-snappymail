<script setup lang="ts">
import { useMailStore } from '../stores/mail'

const mail = useMailStore()
</script>

<template>
  <section class="bg-panel overflow-y-auto min-h-0 flex flex-col">
    <template v-if="mail.selectedMessage">
      <header class="px-4 py-3 border-b border-line">
        <h1 class="text-base font-semibold">{{ mail.selectedMessage.subject }}</h1>
        <p class="text-sm text-ink-sub mt-1">
          From: {{ mail.selectedMessage.from || mail.selectedMessage.fromEmail }}
        </p>
        <p class="text-xs text-ink-mute">{{ mail.selectedMessage.date }}</p>
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
      Select a message
    </div>
  </section>
</template>
