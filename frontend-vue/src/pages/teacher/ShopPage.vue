<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { apiGet, apiPost, apiDelete } from '@/utils/api'
import { useToastStore } from '@/stores/toast'
import type { ApiResponse, ShopItem } from '@/types'

const toast = useToastStore()

const items = ref<ShopItem[]>([])
const loading = ref(true)

onMounted(async () => {
  