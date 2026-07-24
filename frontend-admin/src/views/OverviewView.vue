<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { AdminAPI } from '../api/admin'
import type { Overview } from '../api/types'
import { apiError } from '../api/client'

const ov = ref<Overview | null>(null)
const error = ref('')

function fmt(n: number | null): string {
  return n === null ? 'n/a' : String(n)
}

onMounted(async () => {
  try {
    ov.value = await AdminAPI.overview()
  } catch (e) {
    error.value = apiError(e)
  }
})
</script>

<template>
  <div class="overview">
    <div class="ov-head">Home</div>
    <p v-if="error" class="empty">{{ error }}</p>
    <div v-else-if="ov" class="cards">
      <div class="card"><span class="num">{{ ov.accounts }}</span><span class="lbl">Accounts</span></div>
      <div class="card"><span class="num">{{ ov.domains }}</span><span class="lbl">Domains</span></div>
      <div class="card"><span class="num">{{ ov.aliases }}</span><span class="lbl">Aliases</span></div>
      <div class="card"><span class="num">{{ ov.admins }}</span><span class="lbl">Administrators</span></div>
      <div class="card muted"><span class="num">{{ fmt(ov.servers) }}</span><span class="lbl">Servers</span></div>
      <div class="card muted"><span class="num">{{ fmt(ov.queue) }}</span><span class="lbl">Mail Queue</span></div>
    </div>
  </div>
</template>

<style scoped>
.overview {
  padding: 0;
}
.ov-head {
  padding: 5px 10px;
  background: linear-gradient(to bottom, #fbfbfb, #ededed);
  border-bottom: 1px solid var(--col-header-border);
}
.cards {
  display: flex;
  flex-wrap: wrap;
  gap: 14px;
  padding: 18px;
}
.card {
  width: 150px;
  padding: 16px;
  background: var(--panel);
  border: 1px solid var(--alt);
  border-radius: var(--radius);
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.card .num {
  font-size: 28px;
  color: var(--accent);
}
.card.muted .num {
  color: var(--txt-muted);
}
.card .lbl {
  color: var(--txt);
}
</style>
