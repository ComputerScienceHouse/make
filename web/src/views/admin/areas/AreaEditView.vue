<script setup lang="ts">
import type { Area } from '@/models/areas'
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import DynamicForm, { type FormOptions } from '@/components/DynamicForm.vue'
import router from '@/router'
import type { Training } from '@/models/trainings'
import DynamicTable from '@/components/DynamicTable.vue'
import type { TableOptions } from '@/components/DynamicTable.vue'
import AddTrainingPopup from '@/components/AddTrainingPopup.vue'
import LoadingScreen from '@/components/LoadingScreen.vue'

const area = ref<Area>()
const trainings = ref<Training[]>()
const loading = ref(true)
const notFound = ref(false)

const route = useRoute()

const showTrainingModal = ref(false)

const formOptions: FormOptions<Area> = {
  fields: {
    id: { hidden: true },
  },
}

const tableOptions: TableOptions<Training> = {
  fields: {
    questions: { hidden: true },
  },
  actions: {
    delete: {
      handler: async (t) => {
        await fetch(`/api/trainings/area/${area.value?.id}`, {
          body: JSON.stringify({ id: t.id }),
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

    const [areaRes, trainingRes] = await Promise.all([
      fetch(`/api/areas/${route.params.id}`),
      fetch(`/api/trainings/area/${route.params.id}`),
    ])

    if (areaRes.status === 404) {
      notFound.value = true
    }

    if (!areaRes.ok || !trainingRes.ok) {
      throw new Error('Failed to fetch data')
    }

    area.value = await areaRes.json()
    trainings.value = await trainingRes.json()
  } catch (err) {
    console.error(err)
  } finally {
    loading.value = false
  }
})

async function saveArea(area: Area) {
  const res = await fetch(`/api/areas/${area.id}`, {
    method: 'PUT',
    body: JSON.stringify(area),
    credentials: 'include',
  })

  if (!res.ok) {
    throw new Error(`Failed to save area: ${res.status}`)
  }

  router.push({ path: `/areas/${area.id}` })

  return
}
</script>

<template>
  <AddTrainingPopup
    v-if="showTrainingModal && area"
    :areaid="area.id"
    @close="showTrainingModal = false"
  >
  </AddTrainingPopup>

  <main class="container py-4" v-if="loading">
    <LoadingScreen></LoadingScreen>
  </main>

  <main class="container" v-else-if="area && trainings">
    <div class="d-flex justify-content-between align-items-center">
      <h1>Editing "{{ area?.name }}"</h1>
    </div>

    <DynamicForm :data="area" :options="formOptions" @submit="saveArea" class="mb-3"></DynamicForm>

    <div class="d-flex justify-content-between align-items-center">
      <h1>Required Trainings:</h1>
      <button type="button" class="btn btn-primary" @click="showTrainingModal = true">
        <i class="bi-plus-lg"></i>
      </button>
    </div>

    <DynamicTable :data="trainings" :options="tableOptions"></DynamicTable>
  </main>
</template>

<style scoped>
main {
  padding: 2rem;
  max-width: 100ch;
}
</style>
