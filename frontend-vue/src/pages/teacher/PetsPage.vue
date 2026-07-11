<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { apiGet } from '@/utils/api'
import { getStageEmoji, getStageName } from '@/utils/constants'
import type { ApiResponse, Pet } from '@/types'

const pets = ref<Pet[]>([])
const loading = ref(true)

onMounted(async () => {
  try {
    const res = await apiGet<ApiRespons