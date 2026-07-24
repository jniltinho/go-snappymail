<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { AdminAPI } from '../api/admin'

const auth = useAuthStore()

interface Node {
  name: string
  label: string
  count?: () => number | null
  superadminOnly?: boolean
}

// Live counts for the nav badges (mirrors ZimbraAdmin's 55 / 2 / 0 badges).
const counts = ref<{ accounts: number; aliases: number; domains: number; admins: number } | null>(
  null,
)

const manage: Node[] = [
  { name: 'accounts', label: 'Accounts', count: () => counts.value?.accounts ?? null },
  { name: 'aliases', label: 'Aliases', count: () => counts.value?.aliases ?? null },
  { name: 'distribution-lists', label: 'Distribution Lists', count: () => 0 },
  { name: 'domains', label: 'Domains', count: () => counts.value?.domains ?? null, superadminOnly: true },
  { name: 'admins', label: 'Administrators', count: () => counts.value?.admins ?? null, superadminOnly: true },
]

onMounted(async () => {
  try {
    const o = await AdminAPI.overview()
    counts.value = { accounts: o.accounts, aliases: o.aliases, domains: o.domains, admins: o.admins }
  } catch {
    /* badges stay empty until the backend is reachable */
  }
})
</script>

<template>
  <nav class="nav">
    <div class="home-split">
      <RouterLink :to="{ name: 'overview' }" class="home-main">Home</RouterLink>
      <button class="home-caret" aria-label="Home menu">&#9662;</button>
    </div>

    <div class="node">
      <span class="node-caret">&#9662;</span>
      <span class="node-ico ico-folder"></span>
      <span class="node-label">Manage</span>
    </div>

    <ul class="leaves">
      <template v-for="n in manage" :key="n.name">
        <li v-if="!n.superadminOnly || auth.superadmin">
          <RouterLink :to="{ name: n.name }" active-class="active">
            <span class="leaf-label">{{ n.label }}</span>
            <span v-if="n.count && n.count() !== null" class="badge">{{ n.count() }}</span>
          </RouterLink>
        </li>
      </template>
    </ul>
  </nav>
</template>

<style scoped>
.nav {
  width: 228px;
  flex-shrink: 0;
  background: var(--app-bg);
  border-right: 1px solid var(--panel-border);
  padding: 8px 6px;
  overflow-y: auto;
}

/* "Home" split button */
.home-split {
  display: flex;
  margin-bottom: 8px;
}
.home-main {
  flex: 1;
  padding: 4px 8px;
  background: linear-gradient(to bottom, #ffffff, #ececec);
  border: 1px solid #b7c3cc;
  border-right: none;
  border-radius: var(--radius) 0 0 var(--radius);
  color: var(--txt);
  font-weight: 600;
}
.home-main:hover {
  text-decoration: none;
}
.home-caret {
  width: 20px;
  background: linear-gradient(to bottom, #ffffff, #ececec);
  border: 1px solid #b7c3cc;
  border-radius: 0 var(--radius) var(--radius) 0;
  color: var(--txt);
  font-size: 9px;
  cursor: pointer;
}

/* "Manage" tree node — highlighted selected parent row, like the reference */
.node {
  display: flex;
  align-items: center;
  gap: 5px;
  padding: 4px 4px;
  font-weight: 700;
  color: var(--txt);
  background: linear-gradient(to bottom, #cfe0f1, #b6d2ea);
  border: 1px solid #a7c4e0;
  border-radius: var(--radius);
}
.node-caret {
  font-size: 9px;
  color: #5a6b78;
}
.node-ico {
  width: 16px;
  height: 16px;
  background-size: contain;
  background-repeat: no-repeat;
  background-position: center;
}
.ico-folder {
  background-image: url("data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 16 16'><path fill='%234a7aa8' d='M1 3h5l1 1h8v9H1z'/><path fill='%236699cc' d='M1 5h14v8H1z'/></svg>");
}

/* Leaf rows */
.leaves {
  list-style: none;
  margin: 2px 0 0;
  padding: 0;
}
.leaves a {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 4px 8px 4px 27px;
  color: var(--txt);
}
.leaves a:hover {
  background: #dfe4ea;
  text-decoration: none;
}
.leaves a.active {
  background: var(--sel);
  text-decoration: none;
}
.leaf-label {
  flex: 1;
}
</style>
