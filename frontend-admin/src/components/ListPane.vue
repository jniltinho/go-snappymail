<script setup lang="ts">
// Content-pane wrapper matching ZimbraAdmin: a breadcrumb header bar, a
// toolbar row (New / Edit / Delete), and the list body. Views fill the slots.
defineProps<{ crumb: string }>()
</script>

<template>
  <div class="pane">
    <div class="pane-head">
      <span class="crumb">Home &nbsp;&ndash;&nbsp; Manage &nbsp;&ndash;&nbsp; {{ crumb }}</span>
      <span class="pane-head-right"><a href="#">Help</a></span>
    </div>
    <div class="toolbar">
      <slot name="toolbar" />
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
}
.pane-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 5px 10px;
  background: linear-gradient(to bottom, #fbfbfb, #ededed);
  border-bottom: 1px solid var(--col-header-border);
  font-size: var(--fs);
}
.crumb {
  color: var(--txt);
}
.toolbar {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
  background: linear-gradient(to bottom, #f6f6f6, #e8e8e8);
  border-bottom: 1px solid var(--col-header-border);
  min-height: 30px;
}
.pane-body {
  flex: 1;
  overflow: auto;
  background: var(--view-bg);
}
</style>
