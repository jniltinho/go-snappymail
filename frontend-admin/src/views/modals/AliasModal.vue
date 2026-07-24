<script setup lang="ts">
import { ref } from 'vue'
import BaseModal from '../../components/BaseModal.vue'
import { AdminAPI } from '../../api/admin'
import { apiError } from '../../api/client'
import type { Alias } from '../../api/types'

// One file, both create and edit: `existing` null → New, set → Edit.
const props = defineProps<{ existing: Alias | null }>()
const emit = defineEmits<{ close: []; saved: [] }>()

const address = ref(props.existing?.address ?? '')
const target = ref(props.existing?.goto ?? '')
const active = ref(props.existing?.active ?? true)
const error = ref('')
const busy = ref(false)

async function save() {
  error.value = ''
  busy.value = true
  try {
    const body = { goto: target.value.trim(), active: active.value }
    if (props.existing) {
      await AdminAPI.updateAlias(props.existing.address, body)
    } else {
      await AdminAPI.createAlias({ address: address.value.trim(), ...body })
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
  <BaseModal :title="existing ? `Edit Alias — ${existing.address}` : 'New Alias'" :busy="busy" @close="emit('close')" @submit="save">
    <div class="field">
      <label for="al-a">Alias address:</label>
      <input id="al-a" v-model="address" type="text" :disabled="!!existing" placeholder="alias@example.com" />
    </div>
    <div class="field">
      <label for="al-t">Forward to:</label>
      <input id="al-t" v-model="target" type="text" placeholder="dest1@x, dest2@x" />
    </div>
    <div class="field">
      <label>Status:</label>
      <span class="field-check">
        <input id="al-s" v-model="active" type="checkbox" /><label for="al-s">Active</label>
      </span>
    </div>
    <p v-if="error" class="form-error">{{ error }}</p>
  </BaseModal>
</template>
