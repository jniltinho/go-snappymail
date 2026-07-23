<script setup lang="ts">
import { useMailStore } from '../stores/mail'

const mail = useMailStore()
</script>

<template>
  <section class="border-r border-line bg-panel overflow-y-auto min-h-0">
    <div
      v-for="msg in mail.messages"
      :key="msg.uid"
      class="msg-row"
      :class="{ selected: mail.selectedUid === msg.uid, unread: !msg.seen }"
      @click="mail.selectMessage(msg.uid)"
    >
      <div>
        <div class="text-sm truncate">{{ msg.from || msg.fromEmail }}</div>
        <div class="msg-subject text-sm truncate">{{ msg.subject }}</div>
      </div>
      <div class="text-xs text-ink-mute whitespace-nowrap">{{ msg.date }}</div>
    </div>
    <p v-if="!mail.messages.length" class="p-4 text-sm text-ink-mute">No messages</p>
  </section>
</template>
