<script setup lang="ts" generic="T extends Record<string, any>">
import { computed, reactive, toRaw } from 'vue'

export interface FormOptions<T> {
  fields?: Partial<{
    [K in keyof T]: {
      label?: string
      hidden?: boolean
      type?: 'text' | 'number' | 'checkbox' | 'textarea'
    }
  }>
}

interface Props<T> {
  data: T
  options: FormOptions<T>
}

const props = defineProps<Props<T>>()

const data = reactive({ ...props.data })

// Gets all keys from data expect those hidden from the form options
const keys = computed(() => {
  return (Object.keys(data) as (keyof T)[]).filter((key) => !props.options?.fields?.[key]?.hidden)
})

const emit = defineEmits<{
  submit: [data: T]
}>()

function submit() {
  emit('submit', toRaw(data) as T)
}

function getInputType(key: keyof T): string {
  const override = props.options?.fields?.[key]?.type
  if (override) return override

  const value = (data as T)[key]

  switch (typeof value) {
    case 'number':
      return 'number'
    case 'boolean':
      return 'checkbox'
    default:
      return 'text'
  }
}
</script>
<template>
  <form @submit.prevent="submit">
    <div class="mb-3" v-for="(key, i) in keys" :key="i">
      <label :for="i.toString()" class="form-label">{{ key }}</label>
      <textarea v-if="getInputType(key) === 'textarea'" class="form-control" v-model="data[key]" />

      <input
        v-else-if="getInputType(key) !== 'checkbox'"
        class="form-control"
        :type="getInputType(key)"
        v-model="data[key]"
      />

      <input v-else class="form-check-input" type="checkbox" v-model="data[key]" />
    </div>
    <button type="submit" class="btn btn-primary">Save</button>
  </form>
</template>

<style scoped>
label::first-letter {
  text-transform: capitalize;
}
</style>
