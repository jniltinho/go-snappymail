<script setup lang="ts">
// Base modal shell used by every CRUD dialog (one modal per file). Matches the
// ZimbraAdmin dialog chrome: blue title bar, white body, right-aligned actions.
defineProps<{ title: string; busy?: boolean }>()
const emit = defineEmits<{ close: []; submit: [] }>()
</script>

<template>
  <div class="modal-backdrop" @click.self="emit('close')">
    <div class="modal" role="dialog" :aria-label="title">
      <div class="modal-title">
        <span>{{ title }}</span>
        <button class="modal-x" aria-label="Close" @click="emit('close')">&times;</button>
      </div>
      <form class="modal-body" @submit.prevent="emit('submit')">
        <slot />
        <div class="modal-actions">
          <slot name="actions">
            <button type="button" class="btn" @click="emit('close')">Cancel</button>
            <button type="submit" class="btn btn-primary" :disabled="busy">OK</button>
          </slot>
        </div>
      </form>
    </div>
  </div>
</template>

<style scoped>
.modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.35);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
}
.modal {
  width: 440px;
  max-width: calc(100vw - 24px);
  background: var(--panel);
  border: 1px solid var(--panel-border);
  border-radius: var(--radius);
  box-shadow: 0 4px 18px rgba(0, 0, 0, 0.35);
  overflow: hidden;
}
.modal-title {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 6px 10px;
  color: #fff;
  font-size: var(--fs-big);
  background: linear-gradient(to bottom, var(--dialog-title-top), var(--dialog-title-bottom));
}
.modal-x {
  background: none;
  border: none;
  color: #fff;
  font-size: 16px;
  line-height: 1;
  cursor: pointer;
}
.modal-body {
  padding: 14px 16px;
}
.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-top: 16px;
  padding-top: 12px;
  border-top: 1px solid var(--grid-line);
}
</style>
