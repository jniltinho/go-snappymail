<script setup lang="ts">
import { useMailStore } from '../stores/mail'

const mail = useMailStore()
</script>

<template>
  <div v-if="mail.composeOpen" class="composer-overlay" @keydown.esc="mail.composeOpen = false">
    <div class="composer-window flex flex-col">
      <div class="actionbar flex items-center gap-2 px-3 py-2 border-b border-line">
        <button type="button" class="tbtn" :disabled="mail.composeBusy" @click="mail.sendCompose">
          Send
        </button>
        <button type="button" class="tbtn" :disabled="mail.composeBusy" @click="mail.composeOpen = false">
          Cancel
        </button>
        <button type="button" class="tbtn" :disabled="mail.composeBusy" @click="mail.saveDraftCompose">
          Save Draft
        </button>
        <span v-if="mail.composeErr" class="text-sm login-error px-2 py-0.5 ml-2">{{ mail.composeErr }}</span>
      </div>

      <div class="flex flex-col gap-1 px-3 py-2 border-b border-line">
        <label class="compose-row">
          <span class="compose-label">To:</span>
          <input v-model="mail.cTo" type="text" class="compose-input" placeholder="recipient@example.com" />
        </label>
        <label class="compose-row">
          <span class="compose-label">Cc:</span>
          <input v-model="mail.cCc" type="text" class="compose-input" />
        </label>
        <label class="compose-row">
          <span class="compose-label">Subject:</span>
          <input v-model="mail.cSubject" type="text" class="compose-input" />
        </label>
      </div>

      <textarea v-model="mail.cBody" class="compose-body flex-1 p-3" spellcheck="false"></textarea>
    </div>
  </div>
</template>
