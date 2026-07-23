<script setup lang="ts">
import { onMounted, onBeforeUnmount, watch } from 'vue'
import { useAuthStore } from './stores/auth'
import { useMailStore } from './stores/mail'
import { useSettingsStore } from './stores/settings'
import LoginView from './components/LoginView.vue'
import FolderSidebar from './components/FolderSidebar.vue'
import MessageList from './components/MessageList.vue'
import ReadingPane from './components/ReadingPane.vue'
import AppToolbar from './components/AppToolbar.vue'

const auth = useAuthStore()
const mail = useMailStore()
const settings = useSettingsStore()

function onKey(e: KeyboardEvent) {
  if (['INPUT', 'TEXTAREA'].includes((e.target as HTMLElement).tagName)) return
  const list = mail.messages
  const idx = list.findIndex((m) => m.uid === mail.selectedUid)
  if (e.key === 'j' && idx < list.length - 1) mail.selectMessage(list[idx + 1].uid)
  if (e.key === 'k' && idx > 0) mail.selectMessage(list[idx - 1].uid)
}

watch(
  () => auth.isAuthenticated,
  async (authed) => {
    if (authed) await mail.loadMailbox()
  },
)

onMounted(async () => {
  window.addEventListener('keydown', onKey)
  await auth.checkSession()
  if (auth.isAuthenticated) await mail.loadMailbox()
})

onBeforeUnmount(() => {
  window.removeEventListener('keydown', onKey)
})
</script>

<template>
  <LoginView v-if="!auth.isAuthenticated" />

  <div v-else-if="mail.loading" class="h-full grid place-items-center text-ink-mute">
    Loading mailbox…
  </div>

  <div v-else class="h-full flex flex-col">
    <p v-if="!settings.skinReady(settings.skin)" class="skin-preview-banner">
      Skin preview: {{ settings.skinLabel(settings.skin) }} — full layout coming soon (server: config.toml → ui.skin)
    </p>
    <AppToolbar />
    <div class="grid flex-1 min-h-0 bg-app-bg" style="grid-template-columns: 220px 340px 1fr">
      <FolderSidebar />
      <MessageList />
      <ReadingPane />
    </div>
  </div>
</template>
