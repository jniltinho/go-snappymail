<script setup lang="ts">
import { computed, ref } from 'vue'
import MiniCalendar from './MiniCalendar.vue'

const anchor = ref(new Date())
const mode = ref<'Day' | 'Work Week' | 'Week' | 'Month' | 'List'>('Work Week')
const modes = ['Day', 'Work Week', 'Week', 'Month', 'List'] as const

const days = computed(() => {
  const base = new Date(anchor.value)
  const dow = base.getDay()
  const monday = new Date(base)
  monday.setDate(base.getDate() - ((dow + 6) % 7))
  const count = mode.value === 'Day' ? 1 : mode.value === 'Week' ? 7 : 5
  const start = mode.value === 'Day' ? new Date(anchor.value) : monday
  return Array.from({ length: count }, (_, i) => {
    const d = new Date(start)
    d.setDate(start.getDate() + i)
    return d
  })
})

const rangeLabel = computed(() => {
  const f = days.value[0]
  const l = days.value[days.value.length - 1]
  return `${f.getMonth() + 1}/${f.getDate()} - ${l.getMonth() + 1}/${l.getDate()}`
})

const hours = Array.from({ length: 20 }, (_, i) => i + 4) // 4 AM .. 11 PM

function hourLabel(h: number) {
  if (h === 12) return 'Noon'
  return h < 12 ? `${h} AM` : `${h - 12} PM`
}

function shiftDays(n: number) {
  const d = new Date(anchor.value)
  d.setDate(d.getDate() + n)
  anchor.value = d
}

function isToday(d: Date) {
  return d.toDateString() === new Date().toDateString()
}

function dayHeader(d: Date) {
  return d.toLocaleDateString('en-US', { weekday: 'short', month: 'short', day: 'numeric' })
}
</script>

<template>
  <div class="flex flex-col flex-1 min-h-0">
    <div class="actionbar flex items-center gap-2 px-3 py-1 border-b border-line">
      <div class="new-btn-slot shrink-0">
        <button type="button" class="btn-new w-full text-left px-3" disabled title="Coming soon">
          New Appointment
        </button>
      </div>
      <button type="button" class="tbtn" @click="anchor = new Date()">Today</button>
      <span class="mx-2 flex items-center gap-1 text-sm">
        <button type="button" class="tbtn" @click="shiftDays(mode === 'Day' ? -1 : -7)">◀</button>
        {{ rangeLabel }}
        <button type="button" class="tbtn" @click="shiftDays(mode === 'Day' ? 1 : 7)">▶</button>
      </span>
      <div class="ml-auto flex items-center gap-1">
        <button
          v-for="m in modes"
          :key="m"
          type="button"
          class="tbtn"
          :class="{ 'cal-mode-active': mode === m }"
          @click="mode = m"
        >
          {{ m }}
        </button>
      </div>
    </div>

    <div class="grid flex-1 min-h-0" style="grid-template-columns: 190px 1fr">
      <aside class="border-r border-line bg-panel overflow-y-auto flex flex-col">
        <div class="side-header px-3 py-2">▼ Calendars</div>
        <button type="button" class="side-item w-full text-left active"><span>☑</span><span>Calendar</span></button>
        <button type="button" class="side-item w-full text-left"><span>☐</span><span>Trash</span></button>
        <div class="side-section mt-2"><span>Searches</span><span class="side-gear">⚙</span></div>
        <div class="side-section"><span>Tags</span><span class="side-gear">⚙</span></div>
        <div class="side-section"><span>▸ Zimlets</span></div>
        <MiniCalendar />
      </aside>

      <section class="bg-panel min-h-0 flex flex-col overflow-hidden">
        <template v-if="mode !== 'List' && mode !== 'Month'">
          <div class="cal-headrow" :style="{ gridTemplateColumns: `56px repeat(${days.length}, 1fr)` }">
            <div></div>
            <div v-for="d in days" :key="d.toISOString()" class="cal-dayhead" :class="{ today: isToday(d) }">
              {{ dayHeader(d) }}
            </div>
          </div>
          <div class="flex-1 overflow-y-auto">
            <div
              v-for="h in hours"
              :key="h"
              class="cal-row"
              :style="{ gridTemplateColumns: `56px repeat(${days.length}, 1fr)` }"
            >
              <div class="cal-hour">{{ hourLabel(h) }}</div>
              <div
                v-for="d in days"
                :key="d.toISOString() + h"
                class="cal-cell"
                :class="{ off: h < 8 || h >= 17, today: isToday(d) }"
              ></div>
            </div>
          </div>
        </template>
        <div v-else class="flex-1 grid place-items-center text-ink-mute text-sm">
          {{ mode }} view coming soon.
        </div>
      </section>
    </div>
  </div>
</template>
