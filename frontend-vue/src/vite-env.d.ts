/// <reference types="vite/client" />

declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<object, object, unknown>
  export default component
}

// 宠物系列数据 JSON（src/data/pets/<seriesId>.json）：按系列合并的大对象，
// 用宽松声明避免 TS 对大 JSON 做逐字面量类型推断（拖慢 vue-tsc）。
declare module '*.json' {
  const value: { stories?: Record<string, unknown>; profiles?: Record<string, unknown>; traits?: Record<string, unknown> }
  export default value
}
