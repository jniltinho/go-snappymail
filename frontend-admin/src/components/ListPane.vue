<script setup lang="ts">
import { ref } from 'vue'
// Content-pane wrapper matching ZimbraAdmin: a breadcrumb header bar with a
// Help link and a gear (⚙) actions dropdown on the right, then the list body.
// Actions (New/Edit/Delete) live in the gear menu — the legacy has no toolbar row.
defineProps<{ crumb: string }>()
const menuOpen = ref(false)
function closeMenu() {
  menuOpen.value = false
}
</script>

<template>
  <div class="pane" @click="closeMenu">
    <div class="pane-head">
      <span class="crumb">Home - Manage - {{ crumb }}</span>
      <span class="pane-head-right">
        <span class="help-q">?</span><a href="#" @click.prevent>Help</a>
        <span class="gear-wrap" @click.stop>
          <button class="gear" aria-label="Actions" @click="menuOpen = !menuOpen">
            <span class="gear-ico"></span><span class="caret">&#9662;</span>
          </button>
          <ul v-if="menuOpen" class="gear-menu" @click="menuOpen = false">
            <slot name="menu" />
          </ul>
        </span>
      </span>
    </div>
    <div class="pane-body">
      <slot />
    </div>
  </div>
</template>

<style scoped>
.pane {
  display: flex;
  flex-direction: column;
  height: 100%;
  min-width: 0;
}
.pane-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 5px 10px;
  background: linear-gradient(to bottom, #f6f9fb, #e7eef4);
  border-bottom: 1px solid #c9d3dc;
  font-size: var(--fs);
}
.crumb {
  color: var(--txt);
}
.pane-head-right {
  display: inline-flex;
  align-items: center;
  gap: 8px;
}
.help-q {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 14px;
  height: 14px;
  border-radius: 50%;
  background: var(--accent);
  color: #fff;
  font-size: 10px;
  line-height: 1;
}
.gear-wrap {
  position: relative;
  display: inline-flex;
}
.gear {
  display: inline-flex;
  align-items: center;
  gap: 2px;
  padding: 1px 3px;
  border: 1px solid transparent;
  background: none;
  cursor: pointer;
  color: var(--masthead-icon);
}
.gear:hover {
  border-color: #b7c3cc;
  background: #f0f4f8;
  border-radius: var(--radius);
}
.gear-ico {
  width: 15px;
  height: 15px;
  background: url("data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 16 16' fill='%234a6377'><path d='M8 5a3 3 0 100 6 3 3 0 000-6zm6.5 3l1.2-1-1.2-2-1.5.5a5 5 0 00-1-.6L11.5 3h-3l-.5 1.4a5 5 0 00-1 .6L5.5 4.5l-1.2 2 1.2 1a5 5 0 000 .9l-1.2 1 1.2 2 1.5-.5a5 5 0 001 .6L8 15h3l.5-1.4a5 5 0 001-.6l1.5.5 1.2-2-1.2-1a5 5 0 000-.9z'/></svg>") no-repeat center;
  background-size: 15px 15px;
}
.caret {
  font-size: 9px;
}
.gear-menu {
  position: absolute;
  top: 100%;
  right: 0;
  margin: 3px 0 0;
  padding: 4px 0;
  list-style: none;
  min-width: 150px;
  background: #fff;
  border: 1px solid var(--panel-border);
  border-radius: var(--radius);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
  z-index: 60;
}
.pane-body {
  flex: 1;
  overflow: auto;
  background: var(--view-bg);
}
</style>
