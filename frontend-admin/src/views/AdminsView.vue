<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import ListPane from '../components/ListPane.vue'
import AdminModal from './modals/AdminModal.vue'
import { AdminAPI } from '../api/admin'
import { apiError } from '../api/client'
import type { Admin } from '../api/types'

const rows = ref<Admin[]>([])
const selected = ref<Admin | null>(null)
const error = ref('')
const show = ref(false)
const editing = ref<Admin | null>(null)

const canEdit = computed(() => selected.value !== null)

async function load() {
  error.value = ''
  try {
    rows.value = await AdminAPI.listAdmins()
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
  if (!confirm(`Delete administrator ${selected.value.username}?`)) return
  try {
    await AdminAPI.deleteAdmin(selected.value.username)
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
  <ListPane crumb="Administrators">
    <template #menu>
      <li @click="openNew">New</li>
      <li :class="{ disabled: !canEdit }" @click="canEdit && openEdit()">Edit</li>
      <li :class="{ disabled: !canEdit }" @click="canEdit && remove()">Delete</li>
    </template>

    <p v-if="error" class="empty">{{ error }}</p>
    <table v-else class="grid">
      <thead>
        <tr>
          <th>Username</th>
          <th>Role</th>
          <th>Managed Domains</th>
          <th>Status</th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="a in rows"
          :key="a.username"
          :class="{ selected: selected?.username === a.username }"
          @click="selected = a"
          @dblclick="((selected = a), openEdit())"
        >
          <td class="email namecell ico-admin">{{ a.username }}</td>
          <td>{{ a.superadmin ? 'Super administrator' : 'Domain administrator' }}</td>
          <td>{{ a.superadmin ? 'All' : (a.domains ?? []).join(', ') }}</td>
          <td>{{ a.active ? 'Active' : 'Inactive' }}</td>
        </tr>
        <tr v-if="!rows.length">
          <td colspan="4" class="empty">No administrators.</td>
        </tr>
      </tbody>
    </table>

    <AdminModal v-if="show" :existing="editing" @close="show = false" @saved="onSaved" />
  </ListPane>
</template>
