<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { apiGet } from '@/utils/api'
import { getStageEmoji, getStageName } from '@/utils/constants'
import type { ApiResponse, Student } from '@/types'

const students = ref<Student[]>([])
const loading = ref(true)

onMounted(async () => {
  try {
    const res = await apiGet<ApiRes