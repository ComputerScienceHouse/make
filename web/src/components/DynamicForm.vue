<script setup lang="ts" generic="T extends Record<string, any>">
import { computed, reactive, toRaw } from 'vue';

export interface FormOptions<T> {
  fields?: Partial<{
    [K in keyof T]: {
      label?: string
      hidden?: boolean
    }
  }>
}

interface Props<T> {
    data: T;
    options: FormOptions<T>;
}

const props = defineProps<Props<T>>();

const data = reactive({...props.data})

// Gets all keys from data expect those hidden from the form options
const keys = computed(() => {
  return (Object.keys(data) as (keyof T)[])
    .filter((key) => !props.options?.fields?.[key]?.hidden);
});

const emit = defineEmits<{
  submit: [data: T]
}>()

function submit() {
  emit('submit', toRaw(data) as T)
}
</script>
<template>
    <form @submit.prevent="submit">
        <div class="mb-3" v-for="(key, i) in keys" :key="i">
            <label :for="i.toString()" class="form-label">{{ key }}</label>
            <input type="text" class="form-control" v-model="data[key]" :id="i.toString()">
        </div>
        <button type="submit" class="btn btn-primary">Save</button>
    </form>
</template>