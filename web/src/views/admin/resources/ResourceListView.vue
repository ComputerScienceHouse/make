<script setup lang="ts">
import DynamicTable from '@/components/DynamicTable.vue'
import type { TableOptions } from '@/components/DynamicTable.vue'
import { ref, onMounted } from 'vue'
import LoadingScreen from '@/components/LoadingScreen.vue'
import type { Resource } from '@/models/resources'
import { apiFetch } from '@/util/fetch'

const resources = ref<Resource[]>([])
const loading = ref(true)
const notFound = ref(false)

const tableOptions: TableOptions<Resource> = {
  actions: {
    edit: {
      // relative path, will redirect to
      // /admin/areas/id
      path: (Resource: Resource) => `/admin/resources/${Resource.id}`,
    },
    delete: {
      handler: async (resource: Resource) => {
        await apiFetch(`/api/resources/${resource.id}`, {
          method: 'DELETE',
          credentials: 'include',
        })
      },
    },
  },
}

onMounted(async () => {
  try {
    loading.value = true

    const resourceRes = await apiFetch(`/api/resources/`)

    if (resourceRes.status === 404) {
      notFound.value = true
    }

    if (!resourceRes.ok) {
      throw new Error('Failed to fetch data')
    }

    resources.value = await resourceRes.json()
  } catch (err) {
    console.error(err)
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <main class="container py-4" v-if="loading">
    <LoadingScreen></LoadingScreen>
  </main>

  <main class="container" v-else>
    <div class="d-flex justify-content-between align-items-center">
      <h1>Resources:</h1>
      <RouterLink to="/admin/resources/create" type="button" class="btn btn-primary">Create</RouterLink>
    </div>

    <DynamicTable :data="resources" :options="tableOptions"></DynamicTable>
  </main>
</template>

<style>
main {
  padding: 2rem;
  max-width: 100ch;
}
</style>
