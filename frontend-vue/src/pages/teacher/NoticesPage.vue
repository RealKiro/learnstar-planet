<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { apiGet } from '@/utils/api'
import type { ApiResponse, Notice } from '@/types'

const notices = ref<Notice[]>([])
const loading = ref(true)

onMounted(async () => {
  try {
    const res = await apiGet<ApiResponse<Notice[]>>('/api/v1/teacher/notices')
    notices.value = r