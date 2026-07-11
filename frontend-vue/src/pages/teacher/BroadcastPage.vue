<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { apiGet, apiPost } from '@/utils/api'
import { useToastStore } from '@/stores/toast'
import type { ApiResponse } from '@/types'

const toast = useToastStore()

const bcContent = ref('')
const bcType = ref<'banner' | 'popup' | 'fullscreen'>('banner')
const bcVoice = ref(true)
co