<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import ListPane from '../components/ListPane.vue'
import DomainModal from './modals/DomainModal.vue'
import { AdminAPI } from '../api/admin'
import { apiError } from '../api/client'
import type { Domain } from '../api/types'

const rows = ref<Domain[]>([])
const selected = ref<Domain | null>(null)
const error = ref('')
const show = ref(false)
const editing = ref<Domain | null>(null)

const canEdit = computed(() => selected.value !== null)

async function load() {
  error.value = ''
  try {
    rows.value = await AdminAPI.listDomains()
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
  if (!confirm(`Delete domain ${selected.value.domain}? This removes its accounts and aliases.`)) return
  try {
    await AdminAPI.deleteDomain(selected.value.domain)
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
  <ListPane crumb="Domains">
    <template #toolbar>
      <button class="tbtn" @click="openNew">New</button>
      <button class="tbtn" :disabled="!canEdit" @click="openEdit">Edit</button>
      <button class="tbtn" :disabled="!canEdit" @click="remove">Delete</button>
    </template>

    <p v-if="error" class="empty">{{ error }}</p>
    <table v-else class="grid">
      <thead>
        <tr>
          <th>Domain Name</th>
          <th>Description</th>
          <th>Accounts</th>
          <th>Status</th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="d in rows"
          :key="d.domain"
          :class="{ selected: selected?.domain === d.domain }"
          @click="selected = d"
          @dblclick="((selected = d), openEdit())"
        >
          <td class="email namecell">{{ d.domain }}</td>
          <td>{{ d.description }}</td>
          <td>{{ d.mailboxes }}</td>
          <td>{{ d.active ? 'Active' : 'Inactive' }}</td>
        </tr>
        <tr v-if="!rows.length">
          <td colspan="4" class="empty">No domains.</td>
        </tr>
      </tbody>
    </table>

    <DomainModal v-if="show" :existing="editing" @close="show = false" @saved="onSaved" />
  </ListPane>
</template>
