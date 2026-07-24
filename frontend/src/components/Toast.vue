<script setup lang="ts">
import { useMailStore } from '../stores/mail'

const mail = useMailStore()

function onUndo() {
  mail.toastUndo?.()
}
</script>

<template>
  <transition name="toast-fade">
    <div v-if="mail.toastText" class="gsn-toast">
      {{ mail.toastText }}
      <a v-if="mail.toastUndo" class="gsn-toast-undo" @click="onUndo">Undo</a>
    </div>
  </transition>
</template>
