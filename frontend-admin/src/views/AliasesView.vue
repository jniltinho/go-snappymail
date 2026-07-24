<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import ListPane from '../components/ListPane.vue'
import AliasModal from './modals/AliasModal.vue'
import { AdminAPI } from '../api/admin'
import { apiError } from '../api/client'
import type { Alias } from '../api/types'

const rows = ref<Alias[]>([])
const selected = ref<Alias | null>(null)
const error = ref('')
const show = ref(false)
const editing = ref<Alias | null>(null)

const canEdit = computed(() => selected.value !== null)

async function load() {
  error.value = ''
  try {
    rows.value = await AdminAPI.listAliases()
  } catch (e) {
    error.value = apiError(e)
  }
}

function openNew() {
  editing.value = null
  show.value = true
}
function openEdit() {
  editing.value = selected.value
  show.value = true
}
async function remove() {
  if (!selected.value) return
  if (!confirm(`Delete alias ${selected.value.address}?`)) return
  try {
    await AdminAPI.deleteAlias(selected.value.address)
    selected.value = null
    await load()
  } catch (e) {
    error.value = apiError(e)
  }
}
function onSaved() {
  show.value = false
  load()
}

onMounted(load)
</script>

<template>
  <ListPane crumb="Aliases">
    <template #toolbar>
      <button class="tbtn" @click="openNew">New</button>
      <button class="tbtn" :disabled="!canEdit" @click="openEdit">Edit</button>
      <button class="tbtn" :disabled="!canEdit" @click="remove">Delete</button>
    </template>

    <p v-if="error" class="empty">{{ error }}</p>
    <table v-else class="grid">
      <thead>
        <tr>
          <th>Email Address</th>
          <th>Target</th>
          <th>Status</th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="a in rows"
          :key="a.address"
          :class="{ selected: selected?.address === a.address }"
          @click="selected = a"
          @dblclick="((selected = a), openEdit())"
        >
          <td class="email">{{ a.address }}</td>
          <td>{{ a.goto }}</td>
          <td>{{ a.active ? 'Active' : 'Inactive' }}</td>
        </tr>
        <tr v-if="!rows.length">
          <td colspan="3" class="empty">No aliases.</td>
        </tr>
      </tbody>
    </table>

    <AliasModal v-if="show" :existing="editing" @close="show = false" @saved="onSaved" />
  </ListPane>
</template>
