<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import ListPane from '../components/ListPane.vue'
import NewAccountModal from './modals/NewAccountModal.vue'
import EditAccountModal from './modals/EditAccountModal.vue'
import { AdminAPI } from '../api/admin'
import { apiError } from '../api/client'
import type { Mailbox } from '../api/types'

const rows = ref<Mailbox[]>([])
const selected = ref<Mailbox | null>(null)
const error = ref('')
const showNew = ref(false)
const showEdit = ref(false)

const canEdit = computed(() => selected.value !== null)

async function load() {
  error.value = ''
  try {
    rows.value = await AdminAPI.listMailboxes()
  } catch (e) {
    error.value = apiError(e)
  }
}

async function remove() {
  if (!selected.value) return
  if (!confirm(`Delete account ${selected.value.username}?`)) return
  try {
    await AdminAPI.deleteMailbox(selected.value.username)
    selected.value = null
    await load()
  } catch (e) {
    error.value = apiError(e)
  }
}

function onSaved() {
  showNew.value = false
  showEdit.value = false
  load()
}

onMounted(load)
</script>

<template>
  <ListPane crumb="Accounts">
    <template #menu>
      <li @click="showNew = true">New</li>
      <li :class="{ disabled: !canEdit }" @click="canEdit && (showEdit = true)">Edit</li>
      <li :class="{ disabled: !canEdit }" @click="canEdit && remove()">Delete</li>
    </template>

    <p v-if="error" class="empty">{{ error }}</p>
    <table v-else class="grid">
      <thead>
        <tr>
          <th style="width: 29%">Email Address</th>
          <th style="width: 19%">Display Name</th>
          <th style="width: 12%">Status</th>
          <th style="width: 19%">Last Login</th>
          <th style="width: 21%">Description</th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="m in rows"
          :key="m.username"
          :class="{ selected: selected?.username === m.username }"
          @click="selected = m"
          @dblclick="((selected = m), (showEdit = true))"
        >
          <td class="email namecell ico-account">{{ m.username }}</td>
          <td>{{ m.name }}</td>
          <td>{{ m.active ? 'Active' : 'Inactive' }}</td>
          <td class="muted">Never signed in</td>
          <td></td>
        </tr>
        <tr v-if="!rows.length">
          <td colspan="5" class="empty">No accounts.</td>
        </tr>
      </tbody>
    </table>

    <NewAccountModal v-if="showNew" @close="showNew = false" @saved="onSaved" />
    <EditAccountModal
      v-if="showEdit && selected"
      :account="selected"
      @close="showEdit = false"
      @saved="onSaved"
    />
  </ListPane>
</template>
