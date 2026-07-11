<script setup lang="ts">
import { useToastStore } from '@/stores/toast'

const toastStore = useToastStore()
</script>

<template>
  <div class="toast-container">
    <div
      v-for="t in toastStore.toasts"
      :key="t.id"
      :class="['toast', `toast--${t.type}`]"
    >
      {{ t.message }}
    </div>
  </div>
</template>
