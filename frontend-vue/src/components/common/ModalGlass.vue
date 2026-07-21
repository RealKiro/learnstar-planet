<script setup lang="ts">
defineProps<{ visible: boolean }>()
const emit = defineEmits<{ (e: 'update:visible', v: boolean): void }>()
function close() { emit('update:visible', false) }
</script>

<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="visible" class="modal-overlay" @click.self="close">
        <div class="modal-box">
          <slot />
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.35);
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
}
.modal-box {
  background: rgba(255, 255, 255, 0.92);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-radius: 16px;
  padding: 24px 28px;
  max-width: 540px;
  width: 92%;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.15);
  border: 1px solid rgba(255, 255, 255, 0.2);
}
.modal-enter-active { animation: overlayIn 0.25s ease; }
.modal-leave-active { animation: overlayOut 0.2s ease; }
.modal-enter-active .modal-box { animation: modalPop 0.3s cubic-bezier(0.34, 1.56, 0.64, 1); }
.modal-leave-active .modal-box { animation: modalPopOut 0.2s ease forwards; }
@keyframes overlayIn { 0% { opacity: 0; } 100% { opacity: 1; } }
@keyframes overlayOut { 0% { opacity: 1; } 100% { opacity: 0; } }
@keyframes modalPop { 0% { opacity: 0; transform: scale(0.95) translateY(10px); } 100% { opacity: 1; transform: scale(1) translateY(0); } }
@keyframes modalPopOut { 0% { opacity: 1; transform: scale(1); } 100% { opacity: 0; transform: scale(0.95) translateY(10px); } }
</style>
