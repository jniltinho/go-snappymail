<script setup lang="ts">
import { ref } from 'vue'
import BaseModal from '../../components/BaseModal.vue'
import { AdminAPI } from '../../api/admin'
import { apiError } from '../../api/client'
import type { Mailbox } from '../../api/types'

const props = defineProps<{ account: Mailbox }>()
const emit = defineEmits<{ close: []; saved: [] }>()

const name = ref(props.account.name)
const password = ref('')
const quota = ref(props.account.quota)
const active = ref(props.account.active)
const error = ref('')
const busy = ref(false)

async function save() {
  error.value = ''
  busy.value = true
  try {
    const body: Record<string, unknown> = {
      name: name.value.trim(),
      quota: Number(quota.value) || 0,
      active: active.value,
    }
    if (password.value) body.password = password.value
    await AdminAPI.updateMailbox(props.account.username, body)
    emit('saved')
  } catch (e) {
    error.value = apiError(e)
  } finally {
    busy.value = false
  }
}
</script>

<template>
  <BaseModal :title="`Edit Account — ${account.username}`" :busy="busy" @close="emit('close')" @submit="save">
    <div class="field">
      <label>Email address:</label>
      <span>{{ account.username }}</span>
    </div>
    <div class="field">
      <label for="ea-n">Display name:</label>
      <input id="ea-n" v-model="name" type="text" />
    </div>
    <div class="field">
      <label for="ea-p">New password:</label>
      <input id="ea-p" v-model="password" type="password" placeholder="(leave blank to keep)" />
    </div>
    <div class="field">
      <label for="ea-q">Quota (MB):</label>
      <input id="ea-q" v-model.number="quota" type="number" min="0" />
    </div>
    <div class="field">
      <label>Status:</label>
      <span class="field-check">
        <input id="ea-a" v-model="active" type="checkbox" /><label for="ea-a">Active</label>
      </span>
    </div>
    <p v-if="error" class="form-error">{{ error }}</p>
  </BaseModal>
</template>
