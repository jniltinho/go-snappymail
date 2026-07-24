<script setup lang="ts">
import { onMounted } from 'vue'
import { RouterView } from 'vue-router'
import TopBar from '../components/TopBar.vue'
import NavTree from '../components/NavTree.vue'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()

// Re-hydrate identity from the stored token on a hard reload.
onMounted(() => {
  auth.refresh().catch(() => auth.logout())
})
</script>

<template>
  <div class="shell">
    <TopBar />
    <div class="shell-body">
      <NavTree />
      <main class="content">
        <RouterView />
      </main>
      <aside class="rightrail" aria-hidden="true">
        <button class="rr-collapse" aria-label="Collapse">&#171;</button>
        <div class="rr-count"><i class="rr-dot rr-yellow"></i><span>0</span></div>
        <div class="rr-count"><i class="rr-dot rr-blue"></i><span>0</span></div>
        <div class="rr-count"><i class="rr-dot rr-green"></i><span>0</span></div>
      </aside>
    </div>
    <footer class="statusbar"></footer>
  </div>
</template>

<style scoped>
.shell {
  display: flex;
  flex-direction: column;
  height: 100%;
}
.shell-body {
  flex: 1;
  display: flex;
  min-height: 0;
}
.content {
  flex: 1 1 0;
  min-width: 0;
  width: 0;
  background: var(--view-bg);
  overflow: auto;
}
.rightrail {
  width: 20px;
  flex-shrink: 0;
  background: var(--app-bg);
  border-left: 1px solid var(--panel-border);
  display: flex;
  flex-direction: column;
  align-items: center;
  padding-top: 4px;
  gap: 10px;
}
.rr-collapse {
  border: none;
  background: none;
  color: #6b7a86;
  font-size: 11px;
  cursor: pointer;
  margin-bottom: 4px;
}
.rr-count {
  display: flex;
  flex-direction: column;
  align-items: center;
  font-size: 9px;
  color: var(--txt-muted);
  gap: 1px;
}
.rr-dot {
  width: 10px;
  height: 10px;
  border-radius: 2px;
}
.rr-yellow {
  background: #e0a83e;
}
.rr-blue {
  background: #4a90c2;
}
.rr-green {
  background: #5aa845;
}

.statusbar {
  height: 5px;
  background: var(--statusbar);
}
</style>
