import { ref, reactive } from 'vue'
import { useToastStore } from '@/stores/toast'

type Rule = (value: string) => string | null

interface FieldRule {
  field: string
  rules: Rule[]
  value: () => string
}

export function useInlineValidate() {
  const errors = reactive<Record<string, string>>({})
  const fieldRules = ref<FieldRule[]>([])

  function defineField(field: string, value: () => string, rules: Rule[]) {
    fieldRules.value.push({ field, value, rules })
  }

  function validateField(field: string): boolean {
    const fr = fieldRules.value.find(r => r.field === field)
    if (!fr) return true
    const val = fr.value()
    for (const rule of fr.rules) {
      const err = rule(val)
      if (err) { errors[field] = err; return false }
    }
    delete errors[field]
    return true
  }

  function clearError(field: string) {
    delete errors[field]
  }

  function validateAll(): boolean {
    let valid = true
    for (const fr of fieldRules.value) {
      const val = fr.value()
      for (const rule of fr.rules) {
        const err = rule(val)
        if (err) { errors[fr.field] = err; valid = false; break }
      }
    }
    return valid
  }

  function inputStyle(field: string) {
    return errors[field] ? { borderColor: '#f87171', boxShadow: '0 0 0 2px rgba(248,113,113,0.2)' } : {}
  }

  function errorText(field: string) {
    return errors[field] || ''
  }

  function hasError(field: string) {
    return !!errors[field]
  }

  return {
    errors, defineField, validateField, clearError, validateAll,
    inputStyle, errorText, hasError,
  }
}

// 常用校验规则
export const required = (msg: string) => (val: string) => val.trim() ? null : msg
export const minLength = (min: number, msg?: string) => (val: string) => val.length >= min ? null : (msg || `至少 ${min} 个字符`)
export const isEmail = (msg?: string) => (val: string) => !val || /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(val) ? null : (msg || '邮箱格式不正确')
export const isPhone = (msg?: string) => (val: string) => !val || /^1\d{10}$/.test(val) ? null : (msg || '手机号格式不正确')
