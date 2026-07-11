<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { apiGet, apiPost, apiDelete } from '@/utils/api'
import { useToastStore } from '@/stores/toast'
import { avatarGradient, platformLabel, escapeHtml } from '@/utils/constants'
import type { ApiResponse } from '@/types'

const toast = useToastStore()

interface Teacher {
  id: number
