<script setup lang="ts">
import DynamicTable from '@/components/DynamicTable.vue'
import type { TableOptions } from '@/components/DynamicTable.vue'
import { ref, onMounted } from 'vue'
import type { Training } from '@/models/trainings'
import LoadingScreen from '@/components/LoadingScreen.vue'

const trainings = ref<Training[]>([])
const loading = ref(true)

const tableOptions: TableOptions<Training> = {
  actions: {
    edit: {
      path: (training: Training) => `/admin/trainings/${training.id}`,
    },
    delete: {
      handler: async (training: Training) => {
        await fetch(`/api/trainings/${training.id}`, {
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

    const trainingsRes = await fetch(`/api/trainings/`)

    if (!trainingsRes.ok) {
      throw new Error('Failed to fetch data')
    }

    trainings.value = await trainingsRes.json()
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
      <h1>Trainings:</h1>
      <RouterLink to="/admin/trainings/create" class="btn btn-primary">Create</RouterLink>
    </div>

    <DynamicTable :data="trainings" :options="tableOptions"></DynamicTable>
  </main>
</template>

<style>
main {
  padding: 2rem;
  max-width: 100ch;
}
</style>
