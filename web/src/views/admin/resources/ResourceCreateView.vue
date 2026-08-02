<script setup lang="ts">
import { ref } from 'vue'
import DynamicForm, { type FormOptions } from '@/components/DynamicForm.vue'
import type { Resource } from '@/models/resources'
import { apiFetch } from '@/util/fetch'

const resource = ref<Resource>({
  id: 0, // id is ignored by api when creating a new area
  name: '',
  description: '',
  url: '',
})

const formOptions: FormOptions<Resource> = {
  fields: {
    id: { hidden: true },
  },
}

async function saveResource(resource: Resource) {
  const res = await apiFetch(`/api/resources/`, {
    method: 'POST',
    body: JSON.stringify(resource),
    credentials: 'include',
  })

  if (!res.ok) {
    throw new Error(`Failed to save area: ${res.status}`)
  }

  location.reload()

  return
}
</script>

<template>
  <main class="container" v-if="resource">
    <div class="d-flex justify-content-between align-items-center">
      <h1>Creating New Resource</h1>
    </div>

    <DynamicForm
      :data="resource"
      :options="formOptions"
      @submit="saveResource"
      class="mb-3"
    ></DynamicForm>
  </main>
</template>

<style scoped>
main {
  padding: 2rem;
  max-width: 100ch;
}
</style>
