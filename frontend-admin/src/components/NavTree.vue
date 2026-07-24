<script setup lang="ts">
import { RouterLink } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()

interface Node {
  name: string
  label: string
  superadminOnly?: boolean
}

// "Manage" group mirrors the ZimbraAdmin left tree; Domains/Admins are
// superadmin-only (a domain_admin cannot create domains or other admins).
const manage: Node[] = [
  { name: 'accounts', label: 'Accounts' },
  { name: 'aliases', label: 'Aliases' },
  { name: 'distribution-lists', label: 'Distribution Lists' },
  { name: 'domains', label: 'Domains', superadminOnly: true },
  { name: 'admins', label: 'Administrators', superadminOnly: true },
]
</script>

<template>
  <nav class="nav">
    <RouterLink :to="{ name: 'overview' }" class="nav-home" active-class="active">Home</RouterLink>
    <div class="nav-section">Manage</div>
    <ul class="nav-list">
      <template v-for="n in manage" :key="n.name">
        <li v-if="!n.superadminOnly || auth.superadmin">
          <RouterLink :to="{ name: n.name }" active-class="active">{{ n.label }}</RouterLink>
        </li>
      </template>
    </ul>
  </nav>
</template>

<style scoped>
.nav {
  width: 230px;
  flex-shrink: 0;
  background: var(--app-bg);
  border-right: 1px solid var(--panel-border);
  padding: 8px 0;
  overflow-y: auto;
}
.nav-home {
  display: block;
  margin: 0 8px 8px;
  padding: 4px 8px;
  background: linear-gradient(to bottom, #fff, #eee);
  border: 1px solid var(--alt);
  border-radius: var(--radius);
  color: var(--txt);
  font-weight: 600;
}
.nav-section {
  padding: 4px 12px;
  font-weight: 700;
  color: var(--txt);
}
.nav-list {
  list-style: none;
  margin: 0;
  padding: 0;
}
.nav-list a {
  display: block;
  padding: 4px 12px 4px 26px;
  color: var(--txt);
}
.nav-list a:hover {
  background: #dfe3ea;
  text-decoration: none;
}
.nav-list a.active {
  background: var(--sel);
  text-decoration: none;
}
</style>
