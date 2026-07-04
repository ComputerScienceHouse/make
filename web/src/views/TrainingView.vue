<script setup lang="ts">
import { authState } from '@/auth'
import LoadingScreen from '@/components/LoadingScreen.vue'
import type { SubmissionResults, Training, UserTraining } from '@/models/trainings'
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'

const user = authState.user

const training = ref<Training>()
const userTrainings = ref<UserTraining[]>([])
const loading = ref(true)
const notFound = ref(false)
const completedTraining = ref(-1) // -1 not found, otherwise index of array
const submitted = ref(false)
const submissionResults = ref<SubmissionResults>()

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

    completedTraining.value =
      userTrainings.value?.findIndex((m) => m.trainingId === training.value?.id) ?? -1
  } catch (err) {
    console.error(err)
  } finally {
    loading.value = false
  }
})

async function submitTraining() {
  const payload: Record<string, string> = {}
  loading.value = true

  for (const [key, value] of Object.entries(answers.value)) {
    payload[key] = String(value)
  }

  const res = await fetch(`/api/trainings/${training.value?.id}/submissions`, {
    method: 'POST',
    body: JSON.stringify(payload),
    credentials: 'include',
  })

  if (!res.ok) {
    throw new Error(`Failed to submit training: ${res.status}`)
  }

  submitted.value = true
  loading.value = false
  submissionResults.value = await res.json()
}
</script>

<template>
  <main class="container py-4" v-if="loading">
    <LoadingScreen></LoadingScreen>
  </main>

  <main class="container py-5" v-else-if="completedTraining !== -1">
    <div class="text-center">
      <div>
        <h1 class="h2 mb-3">You've completed this training</h1>
        <p class="text-muted mb-4">No need to do it again</p>

        <div>
          <div class="small text-muted mb-1">Expires At</div>
          <div class="fw-semibold">
            {{ new Date(userTrainings[completedTraining]!.expiresAt).toLocaleString() }}
          </div>
        </div>
      </div>
    </div>
  </main>

  <main class="container py-4" v-else-if="training && !submitted">
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

    <p class="text-muted text-center mt-5" v-if="training.questions.length === 0">
      This training has no questions.
    </p>

    <form @submit.prevent="submitTraining">
      <div v-for="q in training.questions" :key="q.id" class="mb-4 p-3 border rounded shadow-sm">
        <label :for="`${q.id}`" class="form-label fs-6">
          {{ q.label }}
          <span v-if="q.required" class="text-danger">*</span>
        </label>

        <textarea
          v-if="q.type === 'textarea'"
          :id="`${q.id}`"
          v-model="answers[q.id] as string"
          class="form-control"
          :required="q.required"
        ></textarea>

        <input
          v-else-if="q.type === 'number'"
          :id="`${q.id}`"
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
              :name="`${q.id}`"
              v-model="answers[q.id]"
              class="form-check-input pretty-radio"
              :required="q.required"
            />
            <label :for="option" class="form-check-label">{{ option }}</label>
          </div>
        </div>

        <div v-else-if="q.type === 'checkbox'" class="form-check">
          <input
            :id="`${q.id}`"
            v-model="answers[q.id]"
            type="checkbox"
            class="form-check-input"
            :required="q.required"
          />
        </div>
      </div>

      <button class="btn btn-primary" v-if="training.questions.length > 0" type="submit">
        Submit
      </button>
    </form>
  </main>

  <main v-else-if="submitted && submissionResults" class="container py-5">
    <div class="card shadow-sm p-5 text-center mx-auto" style="max-width: 700px">
      <div>
        <i
          class="bi display-1 mb-3"
          :class="submissionResults.passed ? 'text-success bi-check-lg' : 'text-danger bi-x-lg'"
        ></i>
      </div>
      <h1>{{ submissionResults.passed ? 'Training Passed' : 'Training Failed' }}</h1>
      <p class="text-muted">
        {{
          submissionResults.passed
            ? 'You have successfully completed this training!'
            : 'You did not acheive a passing score.'
        }}
      </p>

      <div class="progress mb-4 mt-3">
        <div class="progress-bar" role="progressbar" :style="`width: ${submissionResults.grade}%`">
          {{ submissionResults.grade }}%
        </div>
      </div>

      <div class="row text-center mb-4">
        <div class="col">
          <div class="fs-3 fw-bold text-success">
            {{ submissionResults.numCorrect }}
          </div>
          <div class="text-muted small">Correct</div>
        </div>
        <div class="col">
          <div class="fs-3 fw-bold">
            {{ training?.requiredCorrect }}
          </div>
          <div class="text-muted small">Required Correct</div>
        </div>
        <div class="col">
          <div class="fs-3 fw-bold text-danger">
            {{ submissionResults.numIncorrect }}
          </div>
          <div class="text-muted small">Incorrect</div>
        </div>
      </div>

      <RouterLink class="btn btn-primary" to="/"> Return Home </RouterLink>
    </div>
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
