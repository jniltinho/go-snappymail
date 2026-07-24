<script setup lang="ts">
import { ref } from 'vue'
import MiniCalendar from './MiniCalendar.vue'

interface Task {
  id: number
  text: string
  done: boolean
}

const tasks = ref<Task[]>(JSON.parse(localStorage.getItem('gsn_tasks') || '[]'))
const newText = ref('')
const selectedId = ref<number | null>(null)

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
    </div>

    <div class="grid flex-1 min-h-0" style="grid-template-columns: 190px 1fr">
      <aside class="border-r border-line bg-panel overflow-y-auto flex flex-col">
        <div class="side-header px-3 py-2">▼ Task Lists</div>
        <button type="button" class="side-item w-full text-left active"><span>☑</span><span>Tasks</span></button>
        <button type="button" class="side-item w-full text-left"><span>🗑</span><span>Trash</span></button>
        <div class="side-section mt-2"><span>Searches</span><span class="side-gear">⚙</span></div>
        <div class="side-section"><span>Tags</span><span class="side-gear">⚙</span></div>
        <div class="side-section"><span>▸ Zimlets</span></div>
        <MiniCalendar />
      </aside>

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
          <div v-if="!tasks.length" class="p-8 text-center"><span class="empty-box">No results found.</span></div>
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
    </div>
  </div>
</template>
