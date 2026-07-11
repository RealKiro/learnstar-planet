<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { apiGet } from '@/utils/api'
import type { ApiResponse, LeaderboardEntry } from '@/types'

type LbType = 'total' | 'weekly' | 'pet'

const activeTab = ref<LbType>('total')
const entries = ref<LeaderboardEntry[]>([])
const loading = ref(false)

const tabs: Array<{ key: LbT