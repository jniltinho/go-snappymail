<script setup lang="ts">
import { ref } from 'vue'
import BaseModal from '../../components/BaseModal.vue'
import { AdminAPI } from '../../api/admin'
import { apiError } from '../../api/client'
import type { Domain } from '../../api/types'

const props = defineProps<{ existing: Domain | null }>()
const emit = defineEmits<{ close: []; saved: [] }>()

const domain = ref(props.existing?.domain ?? '')
const description = ref(props.existing?.description ?? '')
const maxquota = ref(props.existing?.maxquota ?? 0)
const active = ref(props.existing?.active ?? true)
const error = ref('')
const busy = ref(false)

async function save() {
  error.value = ''
  busy.value = true
  try {
    const body = {
      description: description.value.trim(),
      maxquota: Number(maxquota.value) || 0,
      active: active.value,
    }
    if (props.existing) {
      await AdminAPI.updateDomain(props.existing.domain, body)
    } else {
      await AdminAPI.createDomain({ domain: domain.value.trim(), ...body })
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
  <BaseModal :title="existing ? `Edit Domain — ${existing.domain}` : 'New Domain'" :busy="busy" @close="emit('close')" @submit="save">
    <div class="field">
      <label for="dm-d">Domain name:</label>
      <input id="dm-d" v-model="domain" type="text" :disabled="!!existing" placeholder="example.com" />
    </div>
    <div class="field">
      <label for="dm-desc">Description:</label>
      <input id="dm-desc" v-model="description" type="text" />
    </div>
    <div class="field">
      <label for="dm-q">Max quota (MB):</label>
      <input id="dm-q" v-model.number="maxquota" type="number" min="0" />
    </div>
    <div class="field">
      <label>Status:</label>
      <span class="field-check">
        <input id="dm-s" v-model="active" type="checkbox" /><label for="dm-s">Active</label>
      </span>
    </div>
    <p v-if="error" class="form-error">{{ error }}</p>
  </BaseModal>
</template>
