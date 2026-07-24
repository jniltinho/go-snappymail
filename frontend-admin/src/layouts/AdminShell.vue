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
  flex: 1;
  min-width: 0;
  background: var(--view-bg);
  overflow: auto;
}
.statusbar {
  height: 5px;
  background: linear-gradient(to bottom, var(--masthead-top), var(--masthead-bottom));
}
</style>
