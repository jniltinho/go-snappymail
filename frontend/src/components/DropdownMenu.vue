<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue'

defineProps<{ label: string; btnClass?: string; alignRight?: boolean; split?: boolean }>()
const emit = defineEmits<{ main: [] }>()

const open = ref(false)
const root = ref<HTMLElement | null>(null)

function onDocClick(e: MouseEvent) {
  if (root.value && !root.value.contains(e.target as Node)) open.value = false
}
function onKey(e: KeyboardEvent) {
  if (e.key === 'Escape') open.value = false
}

onMounted(() => {
  document.addEventListener('click', onDocClick)
  document.addEventListener('keydown', onKey)
})
onBeforeUnmount(() => {
  document.removeEventListener('click', onDocClick)
  document.removeEventListener('keydown', onKey)
})
</script>

<template>
  <div ref="root" class="relative" :class="split ? 'block' : 'inline-block'">
    <template v-if="split">
      <span class="btn-split" :class="btnClass">
        <button type="button" class="btn-split-main" @click="emit('main')">{{ label }}</button>
        <button type="button" class="btn-split-arrow" aria-label="More options" @click="open = !open">
          ▾
        </button>
      </span>
    </template>
    <button v-else type="button" :class="btnClass || 'tbtn'" @click="open = !open">
      {{ label }} ▾
    </button>
    <div v-if="open" class="dd-menu" :class="{ 'dd-right': alignRight }" @click="open = false">
      <slot />
    </div>
  </div>
</template>
