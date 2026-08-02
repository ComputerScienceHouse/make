<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import DynamicForm, { type FormOptions } from '@/components/DynamicForm.vue'
import LoadingScreen from '@/components/LoadingScreen.vue'
import type { Resource } from '@/models/resources'
import { apiFetch } from '@/util/fetch'

const resource = ref<Resource>()
const loading = ref(true)
const notFound = ref(false)

const route = useRoute()

const formOptions: FormOptions<Resource> = {
  fields: {
    id: { hidden: true },
  },
}

onMounted(async () => {
  try {
    loading.value = true

    const resourcesRes = await apiFetch(`/api/resources/${route.params.id}`)

    if (resourcesRes.status === 404) {
      notFound.value = true
    }

    if (!resourcesRes.ok) {
      throw new Error('Failed to fetch data')
    }

    resource.value = await resourcesRes.json()
  } catch (err) {
    console.error(err)
  } finally {
    loading.value = false
  }
})

async function saveResource(resource: Resource) {
  const res = await apiFetch(`/api/resources/${resource.id}`, {
    method: 'PUT',
    body: JSON.stringify(resource),
    credentials: 'include',
  })

  if (!res.ok) {
    throw new Error(`Failed to save resource: ${res.status}`)
  }

  location.reload()

  return
}
</script>

<template>
  <main class="container py-4" v-if="loading">
    <LoadingScreen></LoadingScreen>
  </main>

  <main class="container" v-else-if="resource">
    <div class="d-flex justify-content-between align-items-center">
      <h1>Editing Resource "{{ resource.name }}"</h1>
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
