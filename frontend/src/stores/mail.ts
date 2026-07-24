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

  // Zimbra Classic date style: time today, "Jul 22" this year, "7/20/25" older.
  function fmtListDate(ts: number, fallback: string): string {
    if (!ts) return fallback
    const d = new Date(ts * 1000)
    const now = new Date()
    if (d.toDateString() === now.toDateString()) {
      return d.toLocaleTimeString('en-US', { hour: 'numeric', minute: '2-digit' })
    }
    if (d.getFullYear() === now.getFullYear()) {
      return d.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
    }
    return d.toLocaleDateString('en-US', { month: 'numeric', day: 'numeric', year: '2-digit' })
  }

  function fmtFullDate(ts: number, fallback: string): string {
    if (!ts) return fallback
    const d = new Date(ts * 1000)
    // Zimbra format: "July 23, 2026 9:27 PM" (no "at")
    const date = d.toLocaleDateString('en-US', { month: 'long', day: 'numeric', year: 'numeric' })
    const time = d.toLocaleTimeString('en-US', { hour: 'numeric', minute: '2-digit' })
    return `${date} ${time}`
  }

  function mapMessage(raw: Record<string, unknown>): MailMessage {
    const ts = Number(raw.date_ts) || 0
    return {
      uid: Number(raw.uid),
      subject: String(raw.subject || '(No subject)'),
      from: String(raw.from || ''),
      fromEmail: String(raw.from_email || ''),
      date: fmtListDate(ts, String(raw.date || '')),
      dateFull: fmtFullDate(ts, String(raw.date || '')),
      seen: Boolean(raw.seen),
      flagged: Boolean(raw.flagged),
      size: Number(raw.size) || 0,
      to: String(raw.to || ''),
      hasAttachment: Boolean(raw.has_attachment),
    }
  }

  async function loadFolders(): Promise<void> {
    const res = await axios.get(`${API_BASE}/folders`)
    folders.value = (res.data as Record<string, unknown>[]).map(mapFolder)
    if (!folders.value.some((f) => f.name === currentFolder.value)) {
      currentFolder.value = folders.value.find((f) => f.iconType === 'inbox')?.name || 'INBOX'
    }
  }

  const listPage = ref(1)
  const listHasMore = ref(false)
  const loadingMore = ref(false)

  async function loadMessages(): Promise<void> {
    listPage.value = 1
    const res = await axios.get(`${API_BASE}/mail/${encodeURIComponent(currentFolder.value)}`)
    messages.value = (res.data.messages as Record<string, unknown>[]).map(mapMessage)
    const total = folders.value.find((f) => f.name === currentFolder.value)?.messages ?? 0
    listHasMore.value = messages.value.length < total
    // Zimbra behavior: no auto-open — pane stays empty until the user clicks a row
    if (selectedUid.value && !messages.value.some((m) => m.uid === selectedUid.value)) {
      selectedUid.value = null
    }
    if (selectedUid.value) await loadMessageBody(selectedUid.value)
  }

  // Infinite scroll: append the next page until the whole folder is listed.
  async function loadMoreMessages(): Promise<void> {
    if (loadingMore.value || !listHasMore.value || searchQuery.value.trim()) return
    loadingMore.value = true
    try {
      const next = listPage.value + 1
      const res = await axios.get(`${API_BASE}/mail/${encodeURIComponent(currentFolder.value)}`, {
        params: { page: next },
      })
      const fresh = (res.data.messages as Record<string, unknown>[]).map(mapMessage)
      const seen = new Set(messages.value.map((m) => m.uid))
      messages.value.push(...fresh.filter((m) => !seen.has(m.uid)))
      listPage.value = next
      const total = folders.value.find((f) => f.name === currentFolder.value)?.messages ?? 0
      listHasMore.value = fresh.length > 0 && messages.value.length < total
    } finally {
      loadingMore.value = false
    }
  }

  async function loadMessageBody(uid: number): Promise<void> {
    const msg = messages.value.find((m) => m.uid === uid)
    if (!msg || msg.htmlBody !== undefined) return
    const wasSeen = msg.seen

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
    if (!wasSeen) {
      const f = folders.value.find((x) => x.name === currentFolder.value)
      if (f && f.unseen > 0) f.unseen--
    }
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

  // Zimbra "Read More": advance to the next unread message in the list.
  async function readNextUnread(): Promise<void> {
    const start = messages.value.findIndex((m) => m.uid === selectedUid.value)
    const ordered = [...messages.value.slice(start + 1), ...messages.value.slice(0, start + 1)]
    const next = ordered.find((m) => !m.seen)
    if (next) await selectMessage(next.uid)
  }

  async function refresh(): Promise<void> {
    await loadFolders()
    await loadMessages()
  }

  // Background poll: pick up new mail without a reload. Refreshes the list only
  // when the current folder's counts changed; keeps selection and loaded bodies.
  async function pollNewMail(): Promise<void> {
    if (loading.value || composeOpen.value) return
    try {
      const prev = folders.value.find((f) => f.name === currentFolder.value)
      const prevKey = prev ? `${prev.messages}:${prev.unseen}` : ''
      await loadFolders()
      const cur = folders.value.find((f) => f.name === currentFolder.value)
      const curKey = cur ? `${cur.messages}:${cur.unseen}` : ''
      if (curKey === prevKey || searchQuery.value.trim() || listPage.value > 1) return
      const res = await axios.get(`${API_BASE}/mail/${encodeURIComponent(currentFolder.value)}`)
      const old = new Map(messages.value.map((m) => [m.uid, m]))
      messages.value = (res.data.messages as Record<string, unknown>[]).map((raw) => {
        const fresh = mapMessage(raw)
        const o = old.get(fresh.uid)
        return o?.htmlBody !== undefined
          ? { ...fresh, htmlBody: o.htmlBody, plainBody: o.plainBody, attachments: o.attachments }
          : fresh
      })
    } catch {
      /* offline/unauthenticated — try again next tick */
    }
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
    selectedUid.value = null
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

  async function setSeen(uid: number, seen: boolean): Promise<void> {
    const body = new URLSearchParams()
    body.set('flag', 'seen')
    body.set('value', seen ? '1' : '0')
    await axios.post(
      `${API_BASE}/mail/${encodeURIComponent(currentFolder.value)}/${uid}/flag`,
      body,
    )
    const msg = messages.value.find((m) => m.uid === uid)
    if (msg) msg.seen = seen
    await loadFolders()
  }

  async function moveMessage(uid: number, dest: string): Promise<void> {
    const body = new URLSearchParams()
    body.set('dest', dest)
    await axios.post(
      `${API_BASE}/mail/${encodeURIComponent(currentFolder.value)}/${uid}/move`,
      body,
    )
    messages.value = messages.value.filter((m) => m.uid !== uid)
    if (selectedUid.value === uid) selectedUid.value = messages.value[0]?.uid ?? null
    await loadFolders()
  }

  async function moveToSpecial(iconType: string, names: string[], createName: string): Promise<void> {
    if (!selectedUid.value) return
    let dest =
      folders.value.find((f) => f.iconType === iconType)?.name ??
      folders.value.find((f) => names.includes(f.name.toLowerCase()))?.name
    if (!dest) {
      const body = new URLSearchParams()
      body.set('parent', '')
      body.set('name', createName)
      await axios.post(`${API_BASE}/folders`, body)
      dest = createName
    }
    await moveMessage(selectedUid.value, dest)
  }

  async function archiveSelected(): Promise<void> {
    await moveToSpecial('archive', ['archive'], 'Archive')
  }

  async function spamSelected(): Promise<void> {
    await moveToSpecial('junk', ['junk', 'spam'], 'Junk')
  }

  // ── Composer ──────────────────────────────────────────────────
  const composeOpen = ref(false)
  const composeBusy = ref(false)
  const composeErr = ref('')
  const cTo = ref('')
  const cCc = ref('')
  const cSubject = ref('')
  const cBody = ref('')

  function quoteOriginal(msg: MailMessage): string {
    const body = msg.plainBody || ''
    return [
      '',
      '',
      '----- Original Message -----',
      `From: ${msg.from || msg.fromEmail}`,
      `To: ${msg.to || ''}`,
      `Sent: ${msg.date}`,
      `Subject: ${msg.subject}`,
      '',
      body,
    ].join('\n')
  }

  function openCompose(mode: 'new' | 'reply' | 'replyall' | 'forward' = 'new'): void {
    composeErr.value = ''
    const msg = selectedMessage.value
    if (mode === 'new' || !msg) {
      cTo.value = ''
      cCc.value = ''
      cSubject.value = ''
      cBody.value = ''
    } else {
      const reSubject = /^re:/i.test(msg.subject) ? msg.subject : `Re: ${msg.subject}`
      if (mode === 'reply') {
        cTo.value = msg.fromEmail
        cCc.value = ''
        cSubject.value = reSubject
      } else if (mode === 'replyall') {
        cTo.value = msg.fromEmail
        cCc.value = msg.to || ''
        cSubject.value = reSubject
      } else {
        cTo.value = ''
        cCc.value = ''
        cSubject.value = /^fwd:/i.test(msg.subject) ? msg.subject : `Fwd: ${msg.subject}`
      }
      cBody.value = quoteOriginal(msg)
    }
    composeOpen.value = true
  }

  function composeForm(): URLSearchParams {
    const body = new URLSearchParams()
    body.set('to', cTo.value)
    body.set('cc', cCc.value)
    body.set('subject', cSubject.value)
    body.set('body_plain', cBody.value)
    return body
  }

  async function sendCompose(): Promise<void> {
    composeBusy.value = true
    composeErr.value = ''
    try {
      await axios.post(`${API_BASE}/compose/send`, composeForm())
      composeOpen.value = false
    } catch (e: unknown) {
      const err = e as { response?: { data?: { error?: string } } }
      composeErr.value = err.response?.data?.error || 'Failed to send message.'
    } finally {
      composeBusy.value = false
    }
  }

  async function saveDraftCompose(): Promise<void> {
    composeBusy.value = true
    composeErr.value = ''
    try {
      await axios.post(`${API_BASE}/compose/draft`, composeForm())
      composeOpen.value = false
      await loadFolders()
    } catch (e: unknown) {
      const err = e as { response?: { data?: { error?: string } } }
      composeErr.value = err.response?.data?.error || 'Failed to save draft.'
    } finally {
      composeBusy.value = false
    }
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
    listHasMore,
    loadingMore,
    loadMoreMessages,
    readNextUnread,
    refresh,
    pollNewMail,
    search,
    toggleFlag,
    setSeen,
    moveMessage,
    archiveSelected,
    spamSelected,
    deleteSelected,
    composeOpen,
    composeBusy,
    composeErr,
    cTo,
    cCc,
    cSubject,
    cBody,
    openCompose,
    sendCompose,
    saveDraftCompose,
  }
})
