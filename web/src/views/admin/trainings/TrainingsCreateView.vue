<script setup lang="ts">
import TrainingEditForm from '@/components/TrainingEditForm.vue'
import type { TrainingFull } from '@/models/trainings'
import router from '@/router'
import { ref } from 'vue'

const training = ref<TrainingFull>({
  id: 0,
  title: 'New Training',
  description: 'Gets to creating!',
  requiredCorrect: 0,
  questions: [],
})

async function saveTraining(t: TrainingFull) {
  const res = await fetch(`/api/trainings/`, {
    method: 'POST',
    body: JSON.stringify(t),
    credentials: 'include',
  })

  if (!res.ok) {
    throw new Error(`Failed to save training: ${res.status}`)
  }

  router.push({ path: `/admin/trainings/` })

  return
}
</script>

<template>
  <main class="container py-4">
    <div class="d-flex align-items-center justify-content-between">
      <div class="flex-grow-1 me-3">
        <input type="text" v-model="training.title" class="h1 underline-input" />
        <input type="text" v-model="training.description" class="p underline-input text-muted" />
      </div>
      <div class="text-end flex-shrink-0">
        <button class="btn btn-primary" @click="saveTraining(training)">Publish</button>
      </div>
    </div>

    <hr />

   <TrainingEditForm :training="training"></TrainingEditForm>
  </main>
</template>
