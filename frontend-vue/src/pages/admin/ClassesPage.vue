<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { apiGet, apiDelete } from '@/utils/api'
import { useToastStore } from '@/stores/toast'
import type { ApiResponse, ClassRoom } from '@/types'

const toast = useToastStore()
const classes = ref<ClassRoom[]>([])
const loading = ref(true)

onMounted(async () => {
  try {
    const re