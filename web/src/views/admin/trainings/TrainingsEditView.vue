<script setup lang="ts">
import LoadingScreen from '@/components/LoadingScreen.vue'
import TrainingEditForm from '@/components/TrainingEditForm.vue'
import type { TrainingFull } from '@/models/trainings'
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'

const loading = ref(true)
const notFound = ref(false)
const training = ref<TrainingFull>({
  id: 0,
  title: 'New Training',
  description: 'Gets to creating!',
  requiredCorrect: 0,
  questions: [],
})
const route = useRoute()

onMounted(async () => {
  try {
    loading.value = true

    const trainingRes = await fetch(`/api/trainings/${route.params.id}/full`)

    if (trainingRes.status === 404) {
      notFound.value = true
    }

    if (!trainingRes.ok) {
      throw new Error('Failed to fetch data')
    }

    training.value = await trainingRes.json()
  } catch (err) {
    console.error(err)
  } finally {
    loading.value = false
  }
})

async function saveTraining(t: TrainingFull) {
  const res = await fetch(`/api/trainings/${training.value.id}`, {
    method: 'PUT',
    body: JSON.stringify(t),
    credentials: 'include',
  })

  if (!res.ok) {
    throw new Error(`Failed to save training: ${res.status}`)
  }

  location.reload()

  return
}
</script>

<template>
  <main class="container py-4" v-if="loading">
    <LoadingScreen></LoadingScreen>
  </main>

  <main class="container py-4" v-else-if="training">
    <div class="d-flex align-items-center justify-content-between">
      <div class="flex-grow-1 me-3">
        <input type="text" v-model="training.title" class="h1 underline-input w-100" />
        <input type="text" v-model="training.description" class="p underline-input text-muted" />
        <div class="mt-2">
          <label class="form-label small">Required Correct Answers</label>
          <input
            type="number"
            min="0"
            :max="training.questions.length"
            v-model.number="training.requiredCorrect"
            class="form-control"
            style="max-width: 100px"
          />
        </div>
      </div>
      <div class="text-end flex-shrink-0">
        <button class="btn btn-primary" @click="saveTraining(training)">Update</button>
      </div>
    </div>

    <hr />

    <TrainingEditForm :training="training"></TrainingEditForm>
  </main>
</template>
