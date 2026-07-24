<script setup lang="ts">
import { computed, ref } from 'vue'

const view = ref(new Date())
const today = new Date()

const monthLabel = computed(() =>
  view.value.toLocaleDateString('en-US', { month: 'long', year: 'numeric' }),
)

const weeks = computed(() => {
  const y = view.value.getFullYear()
  const m = view.value.getMonth()
  const first = new Date(y, m, 1)
  const start = new Date(first)
  start.setDate(1 - first.getDay())
  const out: { d: number; month: number; date: Date }[][] = []
  const cur = new Date(start)
  for (let w = 0; w < 6; w++) {
    const row = []
    for (let i = 0; i < 7; i++) {
      row.push({ d: cur.getDate(), month: cur.getMonth(), date: new Date(cur) })
      cur.setDate(cur.getDate() + 1)
    }
    out.push(row)
  }
  return out
})

function shift(months: number) {
  const d = new Date(view.value)
  d.setMonth(d.getMonth() + months)
  view.value = d
}

function isToday(d: Date) {
  return d.toDateString() === today.toDateString()
}
</script>

<template>
  <div class="mini-cal">
    <div class="mini-cal-head">
      <button type="button" @click="shift(-12)">«</button>
      <button type="button" @click="shift(-1)">‹</button>
      <span class="mini-cal-title">{{ monthLabel }}</span>
      <button type="button" @click="shift(1)">›</button>
      <button type="button" @click="shift(12)">»</button>
    </div>
    <table class="mini-cal-grid">
      <thead>
        <tr>
          <th v-for="d in ['S', 'M', 'T', 'W', 'T', 'F', 'S']" :key="d + Math.random()">{{ d }}</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(week, wi) in weeks" :key="wi">
          <td
            v-for="cell in week"
            :key="cell.date.toISOString()"
            :class="{ dim: cell.month !== view.getMonth(), today: isToday(cell.date) }"
          >
            {{ cell.d }}
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
