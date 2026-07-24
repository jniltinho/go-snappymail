<script setup lang="ts">
import { onMounted, onBeforeUnmount, ref, watch } from 'vue'
import { useAuthStore } from './stores/auth'
import { useMailStore } from './stores/mail'
import { useSettingsStore } from './stores/settings'
import LoginView from './components/LoginView.vue'
import FolderSidebar from './components/FolderSidebar.vue'
import MessageList from './components/MessageList.vue'
import ReadingPane from './components/ReadingPane.vue'
import AppToolbar from './components/AppToolbar.vue'
import ComposerModal from './components/ComposerModal.vue'
import ContactsView from './components/ContactsView.vue'
import CalendarView from './components/CalendarView.vue'
import TasksView from './components/TasksView.vue'
import PreferencesView from './components/PreferencesView.vue'

const auth = useAuthStore()
const mail = useMailStore()
const settings = useSettingsStore()

const listWidth = ref(Number(localStorage.getItem('gsn_list_w')) || 305)
const sideWidth = ref(Number(localStorage.getItem('gsn_side_w')) || 190)

function startResize(target: typeof listWidth, key: string, min: number, max: number) {
  return (e: MouseEvent) => {
    e.preventDefault()
    const startX = e.clientX
    const startW = target.value
    const move = (ev: MouseEvent) => {
      target.value = Math.min(max, Math.max(min, startW + ev.clientX - startX))
    }
    const up = () => {
      window.removeEventListener('mousemove', move)
      window.removeEventListener('mouseup', up)
      localStorage.setItem(key, String(target.value))
    }
    window.addEventListener('mousemove', move)
    window.addEventListener('mouseup', up)
  }
}

const startListResize = startResize(listWidth, 'gsn_list_w', 220, 680)
const startSideResize = startResize(sideWidth, 'gsn_side_w', 150, 400)

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
    if (authed) {
      await mail.loadMailbox()
    } else {
      mail.selectedUid = null
      mail.messages = []
    }
  },
)

let pollTimer: ReturnType<typeof setInterval> | undefined

onMounted(async () => {
  window.addEventListener('keydown', onKey)
  await auth.checkSession()
  if (auth.isAuthenticated) await mail.loadMailbox()
  pollTimer = setInterval(() => {
    if (auth.isAuthenticated) void mail.pollNewMail()
  }, 15000)
})

onBeforeUnmount(() => {
  window.removeEventListener('keydown', onKey)
  if (pollTimer) clearInterval(pollTimer)
})
</script>

<template>
  <LoginView v-if="!auth.isAuthenticated" />

  <div v-else-if="mail.loading" class="h-full grid place-items-center text-ink-mute">
    Loading mailbox…
  </div>

  <div v-else class="h-full flex flex-col" :style="{ '--side-w': `${sideWidth}px` }">
    <p v-if="!settings.skinReady(settings.skin)" class="skin-preview-banner">
      Skin preview: {{ settings.skinLabel(settings.skin) }} — full layout coming soon (server: config.toml → ui.skin)
    </p>
    <AppToolbar />
    <div
      v-if="settings.activeTab === 'mail'"
      class="grid flex-1 min-h-0 bg-app-bg"
      :style="{ gridTemplateColumns: `${sideWidth}px 6px ${listWidth}px 6px 1fr` }"
    >
      <FolderSidebar />
      <div class="col-sash" title="Drag to resize" @mousedown="startSideResize"></div>
      <MessageList />
      <div class="col-sash" title="Drag to resize" @mousedown="startListResize"></div>
      <ReadingPane />
    </div>
    <ContactsView v-else-if="settings.activeTab === 'contacts'" />
    <CalendarView v-else-if="settings.activeTab === 'calendar'" />
    <TasksView v-else-if="settings.activeTab === 'tasks'" />
    <PreferencesView v-else-if="settings.activeTab === 'preferences'" />
    <ComposerModal />
  </div>
</template>
