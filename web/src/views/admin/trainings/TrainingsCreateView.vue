<script setup lang="ts">
import TrainingEditForm from '@/components/TrainingEditForm.vue'
import type { TrainingFull } from '@/models/trainings'
import router from '@/router'
import { apiFetch } from '@/util/fetch'
import { ref } from 'vue'

const training = ref<TrainingFull>({
  id: 0,
  title: 'New Training',
  description: 'Gets to creating!',
  requiredCorrect: 0,
  showAnswers: false,
  questions: [],
})

async function saveTraining(t: TrainingFull) {
  const res = await apiFetch(`/api/trainings/`, {
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
                <div class="d-flex mt-2 align-items-center">
          <div class="me-4">
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
          <div class="form-check form-switch">
            <input
              type="checkbox"
              id="required"
              class="form-check-input"
              v-model="training.showAnswers"
            />
            <label class="form-check-label small ms-1">Show answers upon completion </label>
          </div>
        </div>
      </div>
      <div class="text-end flex-shrink-0">
        <button class="btn btn-primary" @click="saveTraining(training)">Publish</button>
      </div>
    </div>

    <hr />

    <TrainingEditForm :training="training"></TrainingEditForm>
  </main>
</template>
