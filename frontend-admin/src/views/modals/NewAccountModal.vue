<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { AdminAPI } from '../../api/admin'
import { apiError } from '../../api/client'

const emit = defineEmits<{ close: []; saved: [] }>()

// Wizard tabs mirror the ZimbraAdmin New Account dialog. Only General
// Information is wired to the backend; the rest are shown for parity.
const TABS = [
  'General Information',
  'Contact Information',
  'Aliases',
  'Member Of',
  'Forwarding',
  'Features',
  'Preferences',
  'Themes',
  'Zimlets',
  'Advanced',
]
const activeTab = ref(0)

const domains = ref<string[]>([])
const accountName = ref('')
const domain = ref('')
const firstName = ref('')
const middleInitial = ref('')
const lastName = ref('')
const displayName = ref('')
const displayAuto = ref(true)
const hideUserGal = ref(false)
const hideAliasesGal = ref(false)
const status = ref('active')
const password = ref('')
const error = ref('')
const busy = ref(false)

const autoDisplay = computed(() =>
  [firstName.value, middleInitial.value, lastName.value].filter(Boolean).join(' '),
)
const effectiveDisplay = computed(() => (displayAuto.value ? autoDisplay.value : displayName.value))

onMounted(async () => {
  try {
    const ds = await AdminAPI.listDomains()
    domains.value = ds.map((d) => d.domain)
    if (domains.value.length) domain.value = domains.value[0]
  } catch {
    /* leave the domain free-typed if the list can't load */
  }
})

async function finish() {
  error.value = ''
  if (!accountName.value.trim() || !domain.value) {
    error.value = 'Account name and domain are required.'
    return
  }
  busy.value = true
  try {
    await AdminAPI.createMailbox({
      username: `${accountName.value.trim()}@${domain.value}`,
      name: effectiveDisplay.value.trim(),
      password: password.value,
      quota: 0,
      active: status.value === 'active',
    })
    emit('saved')
  } catch (e) {
    error.value = apiError(e)
  } finally {
    busy.value = false
  }
}
</script>

<template>
  <div class="wiz-backdrop" @click.self="emit('close')">
    <div class="wiz" role="dialog" aria-label="New Account">
      <div class="wiz-title">
        <span>New Account</span>
        <button class="wiz-expand" aria-label="Expand">&raquo;</button>
      </div>

      <div class="wiz-body">
        <ul class="wiz-tabs">
          <li
            v-for="(t, i) in TABS"
            :key="t"
            :class="{ active: i === activeTab, disabled: i !== 0 }"
            @click="i === 0 && (activeTab = i)"
          >
            {{ t }}
          </li>
        </ul>

        <div class="wiz-pane">
          <template v-if="activeTab === 0">
            <div class="section-head">Account Name</div>
            <div class="wrow">
              <label>Account name:<span class="req">*</span></label>
              <span class="acc-name">
                <input v-model="accountName" type="text" class="in-grow" />
                <span class="at">@</span>
                <select v-model="domain">
                  <option v-for="d in domains" :key="d" :value="d">{{ d }}</option>
                </select>
              </span>
            </div>
            <div class="wrow">
              <label>First name:</label>
              <input v-model="firstName" type="text" />
            </div>
            <div class="wrow">
              <label>Middle initial:</label>
              <input v-model="middleInitial" type="text" class="in-mid" maxlength="4" />
            </div>
            <div class="wrow">
              <label>Last name:<span class="req">*</span></label>
              <input v-model="lastName" type="text" />
            </div>
            <div class="wrow">
              <label>Display name:</label>
              <span class="disp">
                <input
                  :value="displayAuto ? autoDisplay : displayName"
                  :disabled="displayAuto"
                  type="text"
                  class="in-grow"
                  @input="displayName = ($event.target as HTMLInputElement).value"
                />
                <label class="inline-check"><input v-model="displayAuto" type="checkbox" />auto</label>
              </span>
            </div>
            <div class="wrow">
              <label></label>
              <label class="inline-check"><input v-model="hideUserGal" type="checkbox" />Hide user in GAL</label>
            </div>
            <div class="wrow">
              <label></label>
              <label class="inline-check"><input v-model="hideAliasesGal" type="checkbox" />Hide aliases in GAL</label>
            </div>

            <div class="section-head">Account Setup</div>
            <div class="wrow">
              <label>Status:</label>
              <select v-model="status">
                <option value="active">Active</option>
                <option value="maintenance">Maintenance</option>
                <option value="locked">Locked</option>
                <option value="closed">Closed</option>
              </select>
            </div>
            <div class="wrow">
              <label>Password:</label>
              <input v-model="password" type="password" />
            </div>

            <p v-if="error" class="wiz-error">{{ error }}</p>
          </template>

          <div v-else class="wiz-placeholder">This section is not configured in this build.</div>
        </div>
      </div>

      <div class="wiz-foot">
        <button class="btn" @click="emit('close')">Help</button>
        <span class="foot-right">
          <button class="btn" @click="emit('close')">Cancel</button>
          <button class="btn" disabled>Previous</button>
          <button class="btn" disabled>Next</button>
          <button class="btn" :disabled="busy" @click="finish">Finish</button>
        </span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.wiz-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.35);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
}
.wiz {
  width: 720px;
  max-width: calc(100vw - 24px);
  background: var(--panel);
  border: 1px solid #7f8fa0;
  border-radius: var(--radius);
  box-shadow: 0 6px 24px rgba(0, 0, 0, 0.4);
  overflow: hidden;
  display: flex;
  flex-direction: column;
}
.wiz-title {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 5px 10px;
  color: #fff;
  font-size: var(--fs-big);
  background: linear-gradient(to bottom, #6a7f92, #4a5f74);
}
.wiz-expand {
  background: none;
  border: none;
  color: #fff;
  cursor: pointer;
  font-size: 13px;
}
.wiz-body {
  display: flex;
  min-height: 360px;
}
.wiz-tabs {
  width: 160px;
  flex-shrink: 0;
  list-style: none;
  margin: 0;
  padding: 8px 0;
  background: #f2f4f7;
  border-right: 1px solid #cdd4da;
}
.wiz-tabs li {
  padding: 5px 12px 5px 20px;
  cursor: pointer;
  color: var(--txt);
}
.wiz-tabs li.active {
  background: linear-gradient(to bottom, #4a90c2, #2f76ab);
  color: #fff;
}
.wiz-tabs li.disabled {
  color: #8a939c;
  cursor: default;
}
.wiz-pane {
  flex: 1;
  padding: 14px 18px;
  overflow-y: auto;
  max-height: 60vh;
}
.section-head {
  font-weight: 700;
  color: var(--txt);
  border-bottom: 1px solid #cfd6dc;
  padding-bottom: 3px;
  margin: 4px 0 12px;
}
.section-head + .section-head,
.section-head:not(:first-child) {
  margin-top: 20px;
}
.wrow {
  display: grid;
  grid-template-columns: 130px 1fr;
  align-items: center;
  gap: 10px;
  margin-bottom: 8px;
}
.wrow > label:first-child {
  text-align: right;
  color: var(--txt);
}
.req {
  color: #c0392b;
  margin-left: 1px;
}
.wrow input[type='text'],
.wrow input[type='password'],
.wrow select {
  height: 22px;
  padding: 1px 5px;
  border: 1px solid #7f9db9;
  border-radius: 2px;
  font: inherit;
  color: var(--txt);
}
.acc-name {
  display: flex;
  align-items: center;
  gap: 6px;
}
.acc-name .in-grow {
  flex: 1;
}
.acc-name .at {
  color: var(--txt);
}
.disp {
  display: flex;
  align-items: center;
  gap: 8px;
}
.disp .in-grow {
  flex: 1;
}
.in-mid {
  width: 48px;
}
.inline-check {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  color: var(--txt);
}
.wiz-placeholder {
  color: var(--txt-muted);
  padding: 20px 4px;
}
.wiz-error {
  color: #c0392b;
  margin: 10px 0 0;
}
.wiz-foot {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  border-top: 1px solid #cdd4da;
  background: #f2f4f7;
}
.foot-right {
  display: flex;
  gap: 6px;
}
</style>
