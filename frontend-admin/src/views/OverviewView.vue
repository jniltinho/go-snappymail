<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { AdminAPI } from '../api/admin'
import type { Overview } from '../api/types'
import { apiError } from '../api/client'

const ov = ref<Overview | null>(null)
const error = ref('')

function fmt(n: number | null | undefined): string {
  return n === null || n === undefined ? 'n/a' : String(n)
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

    <div v-else class="ov-inner">
      <div class="ov-cols">
        <section class="ov-col">
          <h3>Summary</h3>
          <dl>
            <dt>Zimbra Version:</dt>
            <dd>{{ ov?.version ?? 'n/a' }}</dd>
            <dt>Servers:</dt>
            <dd>{{ fmt(ov?.servers) }}</dd>
            <dt>Accounts:</dt>
            <dd>{{ fmt(ov?.accounts) }}</dd>
            <dt>Domains:</dt>
            <dd>{{ fmt(ov?.domains) }}</dd>
            <dt>Aliases:</dt>
            <dd>{{ fmt(ov?.aliases) }}</dd>
          </dl>
        </section>

        <section class="ov-col">
          <h3>Runtime</h3>
          <dl>
            <dt>Service:</dt>
            <dd>n/a</dd>
            <dt>Active Sessions:</dt>
            <dd>{{ fmt(ov?.active_sessions) }}</dd>
            <dt>Mail Queue:</dt>
            <dd>{{ fmt(ov?.queue) }}</dd>
          </dl>
        </section>
      </div>

      <div class="ov-steps">
        <div class="step">
          <div class="step-title"><span class="num">1</span> Get Started <span class="arrow">&rarr;</span></div>
          <ol>
            <li><a href="#">Install certificates</a></li>
            <li><a href="#">Configure the default class of service</a></li>
          </ol>
        </div>
        <div class="step">
          <div class="step-title"><span class="num">2</span> Configure a Domain <span class="arrow">&rarr;</span></div>
          <ol>
            <li><RouterLink to="/domains">Create a domain</RouterLink></li>
            <li><a href="#">Configure the global address list</a></li>
            <li><a href="#">Configure authentication</a></li>
          </ol>
        </div>
        <div class="step">
          <div class="step-title"><span class="num">3</span> Add Accounts <span class="arrow">&rarr;</span></div>
          <ol>
            <li><RouterLink to="/accounts">Add an account</RouterLink></li>
            <li><RouterLink to="/accounts">Manage accounts</RouterLink></li>
            <li><a href="#">Migration and coexistence</a></li>
          </ol>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.overview {
  height: 100%;
}
.ov-head {
  padding: 5px 10px;
  background: linear-gradient(to bottom, #fbfbfb, #ededed);
  border-bottom: 1px solid var(--col-header-border);
}
.ov-inner {
  padding: 22px 28px;
}
.ov-cols {
  display: flex;
  gap: 120px;
}
.ov-col h3 {
  margin: 0 0 12px;
  font-size: 15px;
  font-weight: 700;
  color: #4a6377;
}
.ov-col dl {
  display: grid;
  grid-template-columns: auto auto;
  gap: 6px 18px;
  margin: 0;
}
.ov-col dt {
  color: var(--txt);
}
.ov-col dd {
  margin: 0;
  color: var(--txt);
  font-weight: 600;
}

.ov-steps {
  display: flex;
  gap: 0;
  margin-top: 40px;
  max-width: 940px;
  background: #f6f8fa;
  border: 1px solid #dbe1e6;
  border-radius: var(--radius);
  padding: 22px 24px;
}
.step {
  flex: 1;
  padding: 0 18px;
}
.step + .step {
  border-left: 1px solid #e2e7ec;
}
.step-title {
  font-size: 15px;
  color: #4a6377;
  margin-bottom: 10px;
}
.step-title .num {
  font-size: 17px;
  font-weight: 700;
}
.step-title .arrow {
  color: #9bb0c1;
}
.step ol {
  margin: 0;
  padding-left: 18px;
  color: var(--txt);
}
.step ol li {
  margin-bottom: 5px;
}
</style>
