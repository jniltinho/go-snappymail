<script setup lang="ts">
import { ref } from 'vue'
import BaseModal from '../../components/BaseModal.vue'
import { AdminAPI } from '../../api/admin'
import { apiError } from '../../api/client'
import type { Admin } from '../../api/types'

const props = defineProps<{ existing: Admin | null }>()
const emit = defineEmits<{ close: []; saved: [] }>()

const username = ref(props.existing?.username ?? '')
const password = ref('')
const superadmin = ref(props.existing?.superadmin ?? false)
const domains = ref((props.existing?.domains ?? []).join(', '))
const active = ref(props.existing?.active ?? true)
const error = ref('')
const busy = ref(false)

function domainList(): string[] {
  return domains.value
    .split(',')
    .map((d) => d.trim())
    .filter(Boolean)
}

async function save() {
  error.value = ''
  busy.value = true
  try {
    if (props.existing) {
      const body: Record<string, unknown> = {
        superadmin: superadmin.value,
        active: active.value,
        domains: superadmin.value ? [] : domainList(),
      }
      if (password.value) body.password = password.value
      await AdminAPI.updateAdmin(props.existing.username, body)
    } else {
      await AdminAPI.createAdmin({
        username: username.value.trim(),
        password: password.value,
        superadmin: superadmin.value,
        active: active.value,
        domains: superadmin.value ? [] : domainList(),
      })
    }
    emit('saved')
  } catch (e) {
    error.value = apiError(e)
  } finally {
    busy.value = false
  }
}
</script>

<template>
  <BaseModal :title="existing ? `Edit Administrator — ${existing.username}` : 'New Administrator'" :busy="busy" @close="emit('close')" @submit="save">
    <div class="field">
      <label for="ad-u">Username:</label>
      <input id="ad-u" v-model="username" type="text" :disabled="!!existing" placeholder="admin@example.com" />
    </div>
    <div class="field">
      <label for="ad-p">{{ existing ? 'New password:' : 'Password:' }}</label>
      <input id="ad-p" v-model="password" type="password" :placeholder="existing ? '(leave blank to keep)' : ''" />
    </div>
    <div class="field">
      <label>Role:</label>
      <span class="field-check">
        <input id="ad-super" v-model="superadmin" type="checkbox" /><label for="ad-super">Super administrator</label>
      </span>
    </div>
    <div v-if="!superadmin" class="field">
      <label for="ad-dom">Managed domains:</label>
      <input id="ad-dom" v-model="domains" type="text" placeholder="a.com, b.com" />
    </div>
    <div class="field">
      <label>Status:</label>
      <span class="field-check">
        <input id="ad-a" v-model="active" type="checkbox" /><label for="ad-a">Active</label>
      </span>
    </div>
    <p v-if="error" class="form-error">{{ error }}</p>
  </BaseModal>
</template>
