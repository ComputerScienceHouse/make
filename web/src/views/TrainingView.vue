<script setup lang="ts">
import { authState } from '@/auth'
import type { Training, UserTraining } from '@/models/trainings'
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'

const user = authState.user

const training = ref<Training | null>(null)
const userTrainings = ref<UserTraining[] | null>(null)
const loading = ref(true)
const notFound = ref(false)

const answers = ref<Record<string, string | number | boolean>>({})

const route = useRoute()

onMounted(async () => {
  try {
    loading.value = true

    const [trainingRes, userTrainingsRes] = await Promise.all([
      fetch(`/api/trainings/${route.params.id}`),
      fetch(`/api/user/${user?.uuid}/trainings`),
    ])

    if (trainingRes.status === 404) {
      notFound.value = true
    }

    if (!trainingRes.ok || !userTrainingsRes.ok) {
      throw new Error('Failed to fetch data')
    }

    training.value = await trainingRes.json()
    userTrainings.value = await userTrainingsRes.json()
  } catch (err) {
    console.error(err)
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <main class="container py-4" v-if="loading">
    <h1>Loading...</h1>
  </main>

  <main class="container py-4" v-else-if="training">
    <div class="d-flex align-items-center justify-content-between">
      <div>
        <h1>{{ training.title }}</h1>
        <p class="text-muted mb-0">{{ training.description }}</p>
      </div>
      <div>
        <p class="mb-1">{{ training.questions.length }} questions</p>
        <p class="m-0">{{ Object.keys(answers).length }} answered</p>
      </div>
    </div>

    <hr />

    <form>
      <div v-for="q in training.questions" :key="q.id" class="mb-4 p-3 border rounded shadow-sm">
        <label :for="q.id" class="form-label fs-6">
          {{ q.label }}
          <span v-if="q.required" class="text-danger">*</span>
        </label>

        <textarea
          v-if="q.type === 'textarea'"
          :id="q.id"
          v-model="answers[q.id] as string"
          class="form-control"
        ></textarea>

        <input
          v-else-if="q.type === 'number'"
          :id="q.id"
          v-model.number="answers[q.id]"
          type="number"
          class="form-control"
          :required="q.required"
        />

        <div v-else-if="q.type == 'radio'">
          <div v-for="option in q.options" :key="option" class="form-check">
            <input
              type="radio"
              :id="option"
              :value="option"
              :name="q.id"
              v-model="answers[q.id]"
              class="form-check-input pretty-radio"
            />
            <label :for="option" class="form-check-label">{{ option }}</label>
          </div>
        </div>

        <div v-else-if="q.type === 'checkbox'" class="form-check">
          <input :id="q.id" v-model="answers[q.id]" type="checkbox" class="form-check-input" />
        </div>
      </div>

      <button class="btn btn-primary">Submit</button>
    </form>
    {{ answers }}
  </main>

  <main class="container py-4" v-else-if="notFound">
    <h1>Area Not Found</h1>
    <p class="text-body-secondary">The requested workshop does not exist.</p>
  </main>

  <main class="container py-4" v-else>
    <h1>Something went wrong.</h1>
    <p class="text-body-secondary">Unable to load this area right now.</p>
  </main>
</template>

<style scoped>
.pretty-radio {
  width: 1.3rem;
  height: 1.3rem;
  margin-top: 0.15rem;
  margin-right: 1rem;
  cursor: pointer;

  border: 1px solid var(--bs-secondary);
}
</style>
