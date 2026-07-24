<script setup lang="ts">
import { computed, ref } from 'vue'
import MiniCalendar from './MiniCalendar.vue'
import DropdownMenu from './DropdownMenu.vue'

interface Task {
  id: number
  text: string
  done: boolean
}

const tasks = ref<Task[]>(JSON.parse(localStorage.getItem('gsn_tasks') || '[]'))
const newText = ref('')
const selectedId = ref<number | null>(null)

const selected = computed(() => tasks.value.find((t) => t.id === selectedId.value) ?? null)

function persist() {
  localStorage.setItem('gsn_tasks', JSON.stringify(tasks.value))
}

function addTask() {
  if (!newText.value.trim()) return
  tasks.value.push({ id: Date.now(), text: newText.value.trim(), done: false })
  newText.value = ''
  persist()
}

function completeSelected() {
  const t = tasks.value.find((x) => x.id === selectedId.value)
  if (t) t.done = true
  persist()
}

function deleteSelected() {
  tasks.value = tasks.value.filter((x) => x.id !== selectedId.value)
  selectedId.value = null
  persist()
}
</script>

<template>
  <div class="flex flex-col flex-1 min-h-0">
    <div class="actionbar flex items-center gap-2 px-3 py-1 border-b border-line">
      <div class="new-btn-slot shrink-0">
        <button type="button" class="btn-new w-full text-left px-3" @click="addTask">New Task</button>
      </div>
      <button type="button" class="tbtn" :disabled="!selectedId">Edit</button>
      <button type="button" class="tbtn" :disabled="!selectedId" @click="deleteSelected">Delete</button>
      <button type="button" class="tbtn" :disabled="!selectedId" @click="completeSelected">
        Mark as Completed
      </button>
      <div class="ml-auto flex items-center gap-1">
        <DropdownMenu label="View" btn-class="tbtn" align-right>
          <button type="button" class="dd-item" disabled>All tasks</button>
        </DropdownMenu>
        <button type="button" class="tbtn cal-arrow" disabled>
          <svg viewBox="0 0 12 12" width="10" height="10"><path d="M8.5 1.5L3.5 6l5 4.5z" fill="currentColor" /></svg>
        </button>
        <button type="button" class="tbtn cal-arrow" disabled>
          <svg viewBox="0 0 12 12" width="10" height="10"><path d="M3.5 1.5L8.5 6l-5 4.5z" fill="currentColor" /></svg>
        </button>
      </div>
    </div>

    <div class="grid flex-1 min-h-0" style="grid-template-columns: 190px 6px 440px 6px 1fr">
      <aside class="border-r border-line bg-panel overflow-y-auto flex flex-col">
        <div class="side-header px-3 py-2">▼ Task Lists</div>
        <button type="button" class="side-item w-full text-left active"><span>☑</span><span>Tasks</span></button>
        <button type="button" class="side-item w-full text-left"><span>🗑</span><span>Trash</span></button>
        <div class="side-section mt-2"><span>Searches</span><span class="side-gear">⚙</span></div>
        <div class="side-section"><span>Tags</span><span class="side-gear">⚙</span></div>
        <div class="side-section"><span>▸ Zimlets</span></div>
        <MiniCalendar />
      </aside>
      <div class="col-sash"></div>

      <section class="bg-panel min-h-0 flex flex-col">
        <div class="sort-header"><span>Sorted by Date ⌄</span></div>
        <div class="border-b border-line bg-panel-2 px-3 py-1.5">
          <input
            v-model="newText"
            class="w-full bg-transparent text-sm outline-none"
            placeholder="Click here to add a new Task"
            @keydown.enter="addTask"
          />
        </div>
        <div class="flex-1 overflow-y-auto">
          <p v-if="!tasks.length" class="p-6 text-sm text-ink-mute text-center">No results found.</p>
          <button
            v-for="t in tasks"
            :key="t.id"
            type="button"
            class="msg-row w-full text-left"
            :class="{ selected: selectedId === t.id }"
            @click="selectedId = t.id"
          >
            <div class="text-sm" :class="{ 'line-through text-ink-mute': t.done }">
              {{ t.done ? '☑' : '☐' }} {{ t.text }}
            </div>
          </button>
        </div>
      </section>
      <div class="col-sash"></div>

      <section class="bg-panel min-h-0 flex flex-col">
        <template v-if="selected">
          <div class="subject-bar flex items-center justify-between px-3">
            <h1 class="subject-title truncate">{{ selected.text }}</h1>
            <span class="text-xs text-ink-sub">{{ selected.done ? 'Completed' : 'Not Started' }}</span>
          </div>
          <div class="p-4 text-sm text-ink-sub">
            Status: {{ selected.done ? 'Completed' : 'Not Started' }}
          </div>
        </template>
        <div v-else class="flex-1 grid place-items-center text-ink-mute text-sm">
          To view a task, click on it.
        </div>
      </section>
    </div>
  </div>
</template>
