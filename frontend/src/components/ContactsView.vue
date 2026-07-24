<script setup lang="ts">
import { computed, ref } from 'vue'
import MiniCalendar from './MiniCalendar.vue'
import DropdownMenu from './DropdownMenu.vue'

interface Contact {
  id: number
  name: string
  email: string
}

const contacts = ref<Contact[]>(JSON.parse(localStorage.getItem('gsn_contacts') || '[]'))
const letter = ref('All')
const selectedId = ref<number | null>(null)
const adding = ref(false)
const fName = ref('')
const fEmail = ref('')

const letters = ['All', '123', ...'ABCDEFGHIJKLMNOPQRSTUVWXYZ']

const filtered = computed(() => {
  if (letter.value === 'All') return contacts.value
  if (letter.value === '123') return contacts.value.filter((c) => /^[0-9]/.test(c.name))
  return contacts.value.filter((c) => c.name.toUpperCase().startsWith(letter.value))
})

function persist() {
  localStorage.setItem('gsn_contacts', JSON.stringify(contacts.value))
}

function saveContact() {
  if (!fName.value.trim()) return
  contacts.value.push({ id: Date.now(), name: fName.value.trim(), email: fEmail.value.trim() })
  persist()
  adding.value = false
  fName.value = ''
  fEmail.value = ''
}

function deleteSelected() {
  contacts.value = contacts.value.filter((c) => c.id !== selectedId.value)
  selectedId.value = null
  persist()
}
</script>

<template>
  <div class="flex flex-col flex-1 min-h-0">
    <div class="actionbar flex items-center gap-2 pr-3 py-1 border-b border-line">
      <div class="new-btn-slot shrink-0">
        <DropdownMenu label="New Contact" btn-class="btn-new" split @main="adding = true">
          <button type="button" class="dd-item" @click="adding = true">New Contact</button>
        </DropdownMenu>
      </div>
      <button type="button" class="tbtn" :disabled="!selectedId">Edit</button>
      <button type="button" class="tbtn" :disabled="!selectedId" @click="deleteSelected">Delete</button>
      <DropdownMenu label="🏷" btn-class="tbtn tbtn-icon">
        <button type="button" class="dd-item" disabled>New Tag…</button>
        <button type="button" class="dd-item" disabled>No tags defined</button>
      </DropdownMenu>
      <DropdownMenu label="Actions" btn-class="tbtn">
        <button type="button" class="dd-item" :disabled="!selectedId" @click="deleteSelected">Delete</button>
      </DropdownMenu>
    </div>

    <div class="grid flex-1 min-h-0" style="grid-template-columns: 190px 1fr">
      <aside class="border-r border-line bg-panel overflow-y-auto flex flex-col">
        <div class="side-header px-3 py-2">▼ Contact Lists</div>
        <button type="button" class="side-item w-full text-left active">
          <span>👤</span><span>Contacts</span>
        </button>
        <button type="button" class="side-item w-full text-left"><span>👥</span><span>Distribution Lists</span></button>
        <button type="button" class="side-item w-full text-left"><span>✉</span><span>Emailed Contacts</span></button>
        <button type="button" class="side-item w-full text-left"><span>🗑</span><span>Trash</span></button>
        <div class="side-section mt-2"><span>Searches</span><span class="side-gear">⚙</span></div>
        <div class="side-section"><span>Tags</span><span class="side-gear">⚙</span></div>
        <div class="side-section"><span>▸ Zimlets</span></div>
        <MiniCalendar />
      </aside>

      <section class="bg-panel min-h-0 flex flex-col">
        <div class="alpha-bar">
          <button
            v-for="l in letters"
            :key="l"
            type="button"
            class="alpha-item"
            :class="{ active: letter === l }"
            @click="letter = l"
          >
            {{ l }}
          </button>
        </div>

        <div v-if="adding" class="p-4 border-b border-line flex items-center gap-2">
          <input v-model="fName" placeholder="Full name" class="compose-input max-w-56" />
          <input v-model="fEmail" placeholder="Email" class="compose-input max-w-64" />
          <button type="button" class="tbtn" @click="saveContact">Save</button>
          <button type="button" class="tbtn" @click="adding = false">Cancel</button>
        </div>

        <div class="flex-1 overflow-y-auto p-4">
          <div v-if="!filtered.length" class="empty-box">No results found.</div>
          <button
            v-for="c in filtered"
            :key="c.id"
            type="button"
            class="msg-row w-full text-left"
            :class="{ selected: selectedId === c.id }"
            @click="selectedId = c.id"
          >
            <div>
              <div class="text-sm font-bold">{{ c.name }}</div>
              <div class="text-sm text-ink-sub">{{ c.email }}</div>
            </div>
          </button>
        </div>
      </section>
    </div>
  </div>
</template>
