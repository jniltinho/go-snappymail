import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { axios } from '../api/client'
import type { Folder, MailMessage } from '../types'

export const useMailStore = defineStore('mail', () => {
  const loading = ref(false)
  const folders = ref<Folder[]>([])
  const messages = ref<MailMessage[]>([])
  const currentFolder = ref('INBOX')
  const selectedUid = ref<number | null>(null)
  const searchQuery = ref('')

  const selectedMessage = computed(() =>
    messages.value.find((m) => m.uid === selectedUid.value) ?? null,
  )

  function mapFolder(raw: Record<string, unknown>): Folder {
    return {
      name: String(raw.Name ?? ''),
      label: String(raw.DisplayName || raw.Name || ''),
      iconType: String(raw.IconType || 'folder'),
      unseen: Number(raw.Unseen) || 0,
      messages: Number(raw.Messages) || 0,
      depth: Number(raw.Depth) || 0,
    }
  }

  function mapMessage(raw: Record<string, unknown>): MailMessage {
    return {
      uid: Number(raw.uid),
      subject: String(raw.subject || '(No subject)'),
      from: String(raw.from || ''),
      fromEmail: String(raw.from_email || ''),
      date: String(raw.date || ''),
      seen: Boolean(raw.seen),
      flagged: Boolean(raw.flagged),
      size: Number(raw.size) || 0,
    }
  }

  async function loadFolders(): Promise<void> {
    const res = await axios.get(`${API_BASE}/folders`)
    folders.value = (res.data as Record<string, unknown>[]).map(mapFolder)
    if (!folders.value.some((f) => f.name === currentFolder.value)) {
      currentFolder.value = folders.value.find((f) => f.iconType === 'inbox')?.name || 'INBOX'
    }
  }

  async function loadMessages(): Promise<void> {
    const res = await axios.get(`${API_BASE}/mail/${encodeURIComponent(currentFolder.value)}`)
    messages.value = (res.data.messages as Record<string, unknown>[]).map(mapMessage)
    if (messages.value.length && !messages.value.some((m) => m.uid === selectedUid.value)) {
      selectedUid.value = messages.value[0]?.uid ?? null
    }
    if (selectedUid.value) await loadMessageBody(selectedUid.value)
  }

  async function loadMessageBody(uid: number): Promise<void> {
    const msg = messages.value.find((m) => m.uid === uid)
    if (!msg || msg.htmlBody !== undefined) return

    const res = await axios.get(
      `${API_BASE}/mail/${encodeURIComponent(currentFolder.value)}/${uid}`,
    )
    msg.htmlBody = res.data.html_body || ''
    msg.plainBody = res.data.plain_body || ''
    msg.attachments = (res.data.attachments || []).map((a: Record<string, unknown>) => ({
      filename: String(a.filename || ''),
      part: Number(a.part),
      sizeLabel: String(a.size_label || ''),
      contentType: String(a.content_type || ''),
    }))
    msg.seen = true
  }

  async function selectFolder(name: string): Promise<void> {
    currentFolder.value = name
    selectedUid.value = null
    await loadMessages()
  }

  async function selectMessage(uid: number): Promise<void> {
    selectedUid.value = uid
    await loadMessageBody(uid)
  }

  async function refresh(): Promise<void> {
    await loadFolders()
    await loadMessages()
  }

  async function search(): Promise<void> {
    if (!searchQuery.value.trim()) {
      await loadMessages()
      return
    }
    const res = await axios.get(`${API_BASE}/search`, {
      params: { q: searchQuery.value, mailbox: currentFolder.value },
    })
    messages.value = (res.data.messages as Record<string, unknown>[]).map(mapMessage)
    selectedUid.value = messages.value[0]?.uid ?? null
    if (selectedUid.value) await loadMessageBody(selectedUid.value)
  }

  async function loadMailbox(): Promise<void> {
    loading.value = true
    try {
      await loadFolders()
      await loadMessages()
    } finally {
      loading.value = false
    }
  }

  async function toggleFlag(uid: number, flagged: boolean): Promise<void> {
    const body = new URLSearchParams()
    body.set('flag', 'flagged')
    body.set('value', flagged ? '1' : '0')
    await axios.post(
      `${API_BASE}/mail/${encodeURIComponent(currentFolder.value)}/${uid}/flag`,
      body,
    )
    const msg = messages.value.find((m) => m.uid === uid)
    if (msg) msg.flagged = flagged
  }

  async function deleteSelected(): Promise<void> {
    if (!selectedUid.value) return
    await axios.delete(
      `${API_BASE}/mail/${encodeURIComponent(currentFolder.value)}/${selectedUid.value}`,
    )
    messages.value = messages.value.filter((m) => m.uid !== selectedUid.value)
    selectedUid.value = messages.value[0]?.uid ?? null
  }

  return {
    loading,
    folders,
    messages,
    currentFolder,
    selectedUid,
    selectedMessage,
    searchQuery,
    loadMailbox,
    selectFolder,
    selectMessage,
    refresh,
    search,
    toggleFlag,
    deleteSelected,
  }
})
