<script setup lang="ts">
import { ref } from 'vue'
import BaseModal from '../../components/BaseModal.vue'
import { AdminAPI } from '../../api/admin'
import { apiError } from '../../api/client'

const emit = defineEmits<{ close: []; saved: [] }>()

const username = ref('')
const name = ref('')
const password = ref('')
const quota = ref(0)
const active = ref(true)
const error = ref('')
const busy = ref(false)

async function save() {
  error.value = ''
  busy.value = true
  try {
    await AdminAPI.createMailbox({
      username: username.value.trim(),
      name: name.value.trim(),
      password: password.value,
      quota: Number(quota.value) || 0,
      active: active.value,
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
  <BaseModal title="New Account" :busy="busy" @close="emit('close')" @submit="save">
    <div class="field">
      <label for="na-u">Email address:</label>
      <input id="na-u" v-model="username" type="text" placeholder="user@example.com" />
    </div>
    <div class="field">
      <label for="na-n">Display name:</label>
      <input id="na-n" v-model="name" type="text" />
    </div>
    <div class="field">
      <label for="na-p">Password:</label>
      <input id="na-p" v-model="password" type="password" />
    </div>
    <div class="field">
      <label for="na-q">Quota (MB):</label>
      <input id="na-q" v-model.number="quota" type="number" min="0" />
    </div>
    <div class="field">
      <label>Status:</label>
      <span class="field-check">
        <input id="na-a" v-model="active" type="checkbox" /><label for="na-a">Active</label>
      </span>
    </div>
    <p v-if="error" class="form-error">{{ error }}</p>
  </BaseModal>
</template>
