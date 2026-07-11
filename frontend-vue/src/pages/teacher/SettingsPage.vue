<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { apiGet, apiPost } from '@/utils/api'
import { useAuthStore } from '@/stores/auth'
import { useToastStore } from '@/stores/toast'
import { platformLabel } from '@/utils/constants'
import type { ApiResponse } from '@/types'

const authStore = useAuthStore()
const toast = useTo