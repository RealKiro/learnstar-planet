<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { apiGet, apiPost } from '@/utils/api'
import { useToastStore } from '@/stores/toast'
import type { ApiResponse, ScoreRule, Student } from '@/types'

const toast = useToastStore()

const students = ref<Student[]>([])
const rules = ref<ScoreRule[]>([])
const score