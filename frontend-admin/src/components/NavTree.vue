<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { AdminAPI } from '../api/admin'

const auth = useAuthStore()
const route = useRoute()

interface Leaf {
  name: string
  label: string
  count?: () => number | null
  superadminOnly?: boolean
}

const counts = ref<{ accounts: number; aliases: number; domains: number; admins: number } | null>(
  null,
)

const manage: Leaf[] = [
  { name: 'accounts', label: 'Accounts', count: () => counts.value?.accounts ?? null },
  { name: 'aliases', label: 'Aliases', count: () => counts.value?.aliases ?? null },
  { name: 'distribution-lists', label: 'Distribution Lists', count: () => 0 },
  { name: 'domains', label: 'Domains', count: () => counts.value?.domains ?? null, superadminOnly: true },
  { name: 'admins', label: 'Administrators', count: () => counts.value?.admins ?? null, superadminOnly: true },
]

const MANAGE_ROUTES = new Set(manage.map((m) => m.name))
const onHome = computed(() => route.name === 'overview')
const inManage = computed(() => MANAGE_ROUTES.has(String(route.name)))

// Zimbra sections present in this backend only for Manage; the others mirror
// the legacy tree structure but are inert (not implemented here).
const sections = [
  { key: 'monitor', label: 'Monitor', ico: 'ico-monitor' },
  { key: 'configure', label: 'Configure', ico: 'ico-configure' },
  { key: 'tools', label: 'Tools and Migration', ico: 'ico-tools' },
  { key: 'search', label: 'Search', ico: 'ico-search' },
  { key: 'help', label: 'Help Center', ico: 'ico-help' },
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

    <ul class="tree">
      <!-- Home (leaf) -->
      <li>
        <RouterLink :to="{ name: 'overview' }" class="tnode" :class="{ sel: onHome }">
          <span class="tnode-caret-spacer"></span>
          <span class="tnode-ico ico-home"></span><span class="tnode-label">Home</span>
        </RouterLink>
      </li>

      <!-- Monitor (inert, collapsed) -->
      <li>
        <span class="tnode inert">
          <span class="tnode-caret">&#9656;</span>
          <span class="tnode-ico ico-monitor"></span><span class="tnode-label">Monitor</span>
        </span>
      </li>

      <!-- Manage (expanded section bar) -->
      <li>
        <span class="tnode section" :class="{ sel: inManage }">
          <span class="tnode-caret open">&#9662;</span>
          <span class="tnode-ico ico-manage"></span><span class="tnode-label">Manage</span>
        </span>
      </li>
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

      <!-- Remaining inert sections -->
      <li v-for="s in sections.slice(1)" :key="s.key">
        <span class="tnode inert">
          <span v-if="s.key !== 'help'" class="tnode-caret">&#9656;</span>
          <span v-else class="tnode-caret-spacer"></span>
          <span class="tnode-ico" :class="s.ico"></span><span class="tnode-label">{{ s.label }}</span>
        </span>
      </li>
    </ul>
  </nav>
</template>

<style scoped>
.nav {
  width: 228px;
  flex-shrink: 0;
  background: var(--app-bg);
  border-right: 1px solid var(--panel-border);
  padding: 8px 0;
  overflow-y: auto;
}

/* "Home" split button */
.home-split {
  display: flex;
  margin: 0 6px 8px;
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

.tree {
  list-style: none;
  margin: 0;
  padding: 0;
}

/* Top-level tree node */
.tnode {
  display: flex;
  align-items: center;
  gap: 5px;
  padding: 4px 8px;
  color: var(--txt);
}
.tnode.inert {
  cursor: default;
}
.tnode:hover:not(.inert):not(.sel) {
  background: #dfe4ea;
  text-decoration: none;
}
/* Selected section gets a full-width blue bar */
.tnode.sel,
.tnode.section.sel {
  background: linear-gradient(to bottom, #d5e6f6, #bcd8ef);
  border-top: 1px solid #c6dcf1;
  border-bottom: 1px solid #a9c8e6;
  font-weight: 700;
}
.tnode.section {
  font-weight: 700;
}
.tnode-caret {
  font-size: 8px;
  color: #5a6b78;
  width: 9px;
  text-align: center;
}
.tnode-caret.open {
  font-size: 9px;
}
.tnode-caret-spacer {
  display: inline-block;
  width: 9px;
}
.tnode-ico {
  width: 16px;
  height: 16px;
  background-size: 16px 16px;
  background-repeat: no-repeat;
  background-position: center;
  flex-shrink: 0;
}
.ico-home {
  background-image: url('../assets/icons/ImgHome.png');
}
.ico-monitor {
  background-image: url('../assets/icons/ImgMonitor.png');
}
.ico-manage {
  background-image: url('../assets/icons/ImgManageAccounts.png');
}
.ico-configure {
  background-image: url('../assets/icons/ImgConfigure.png');
}
.ico-tools {
  background-image: url('../assets/icons/ImgToolsAndMigration.png');
}
.ico-search {
  background-image: url('../assets/icons/ImgSearchAll.png');
}
.ico-help {
  background-image: url("data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 16 16'><circle cx='8' cy='8' r='7' fill='%232f7fb5'/><text x='8' y='12' font-size='11' fill='white' text-anchor='middle' font-family='sans-serif'>?</text></svg>");
}

/* The Home leaf node has no caret — align it with the caret'd nodes */
.tnode.sel .tnode-ico,
li > .tnode:not(.section) .tnode-ico {
  margin-left: 0;
}

/* Manage leaf rows */
.leaves {
  list-style: none;
  margin: 0;
  padding: 0;
}
.leaves a {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 4px 8px 4px 40px;
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
