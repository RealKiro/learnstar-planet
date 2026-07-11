import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useToastStore } from '@/stores/toast'

export function createCrudStore<T extends { id: number }>(
  name: string,
  fetchFn: () => Promise<T[]>,
  createFn?: (data: Partial<T>) => Promise<T>,
  updateFn?: (id: number, data: Partial<T>) => Promise<T>,
  deleteFn?: (id: number) => Promise<void>,
) {
  return defineStore(name, () => {
    const items = ref<T[]>([])
    const loading = ref(false)
    const toast = useToastStore()

    async function fetch() {
      loading.value = true
      try {
        items.value = await fetchFn()
      } catch (e) {
        // error handled by interceptor
      } finally {
        loading.value = false
      }
    }

    async function create(data: Partial<T>) {
      if (!createFn) return
      try {
        await createFn(data)
        toast.show('创建成功', 'success')
        await fetch()
      } catch { /* handled */ }
    }

    async function update(id: number, data: Partial<T>) {
      if (!updateFn) return
      try {
        await updateFn(id, data)
        toast.show('更新成功', 'success')
        await fetch()
      } catch { /* handled */ }
    }

    async function remove(id: number) {
      if (!deleteFn) return
      try {
        await deleteFn(id)
        toast.show('删除成功', 'success')
        await fetch()
      } catch { /* handled */ }
    }

    return { items, loading, fetch, create, update, remove }
  })
}
