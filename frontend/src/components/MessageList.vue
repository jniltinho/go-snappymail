<script setup lang="ts">
import { computed } from 'vue'
import { useMailStore } from '../stores/mail'

const mail = useMailStore()

const totalMessages = computed(
  () => mail.folders.find((f) => f.name === mail.currentFolder)?.messages ?? mail.messages.length,
)
</script>

<template>
  <section class="bg-panel min-h-0 flex flex-col">
    <div class="sort-header">
      <span>Sorted by Date</span>
      <span>{{ totalMessages }} message{{ totalMessages === 1 ? '' : 's' }}</span>
    </div>

    <div class="flex-1 overflow-y-auto min-h-0">
      <div
        v-for="msg in mail.messages"
        :key="msg.uid"
        class="msg-row"
        :class="{ selected: mail.selectedUid === msg.uid, unread: !msg.seen }"
        @click="mail.selectMessage(msg.uid)"
      >
        <div class="min-w-0">
          <div class="text-sm truncate">{{ msg.from || msg.fromEmail }}</div>
          <div class="msg-subject text-sm truncate">{{ msg.subject }}</div>
        </div>
        <div class="text-right shrink-0">
          <div class="msg-date text-xs text-ink-sub whitespace-nowrap">{{ msg.date }}</div>
          <span v-if="msg.hasAttachment" class="row-attach" title="Has attachment">📎</span>
          <span
            class="row-flag"
            :class="{ on: msg.flagged }"
            title="Flag"
            @click.stop="mail.toggleFlag(msg.uid, !msg.flagged)"
            >⚑</span
          >
        </div>
      </div>
      <p v-if="!mail.messages.length" class="p-4 text-sm text-ink-mute">No messages</p>
    </div>
  </section>
</template>
