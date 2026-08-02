<script setup lang="ts">
import { authState } from '@/auth'
import LoadingScreen from '@/components/LoadingScreen.vue'
import type { Area } from '@/models/areas'
import type { Training, UserTraining } from '@/models/trainings'
import type { Resource } from '@/models/resources'
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { apiFetch } from '@/util/fetch'

const user = authState.user

const area = ref<Area | null>(null)
const trainings = ref<Training[] | null>(null)
const userTrainings = ref<UserTraining[] | null>(null)
const resources = ref<Resource[] | null>(null)
const loading = ref(true)
const notFound = ref(false)

const completedAllTrainings = computed(() => {
  return (
    trainings.value?.every((t) => userTrainings.value?.some((tr) => tr.trainingId == t.id)) || false
  )
})

const progress = computed(() => {
  if (!trainings.value?.length) return 0

  const requiredTrainings = trainings.value.map((t) => t.id)
  const completedRequiredTrainings = userTrainings.value?.filter((training) =>
    requiredTrainings.includes(training.trainingId),
  )
  return ((completedRequiredTrainings?.length || 0) / requiredTrainings.length) * 100
})

const route = useRoute()
const router = useRouter()

async function deleteArea() {
  try {
    const response = await fetch(`/api/areas/${area.value?.id}`, {
      method: 'DELETE',
      credentials: 'include',
    })

    if (!response.ok) {
      throw new Error(`Failed to delete area (${response.status})`)
    }
  } catch (error) {
    console.error('Error deleting area:', error)
  }

  router.push({ path: '/' })
}

onMounted(async () => {
  try {
    loading.value = true

    const [areaRes, trainingsRes, userTrainingsRes, resourcesRes] = await Promise.all([
      apiFetch(`/api/areas/${route.params.id}`),
      apiFetch(`/api/trainings/area/${route.params.id}`),
      apiFetch(`/api/user/${user?.uuid}/trainings`),
      apiFetch(`/api/resources/area/${route.params.id}`),
    ])

    if (areaRes.status === 404) {
      notFound.value = true
    }

    if (!areaRes.ok || !trainingsRes.ok || !userTrainingsRes.ok || !resourcesRes.ok) {
      throw new Error('Failed to fetch data')
    }

    area.value = await areaRes.json()
    trainings.value = await trainingsRes.json()
    userTrainings.value = await userTrainingsRes.json()
    resources.value = await resourcesRes.json()
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

  <main class="container py-4" v-else-if="area">
    <div class="position-relative rounded overflow-hidden mb-4" style="height: 280px">
      <img :src="area.photourl" :alt="area.name" class="w-100 h-100 object-fit-cover" />

      <div class="position-absolute top-0 start-0 w-100 h-100 image-gradient"></div>

      <!-- admin only button -->
      <div v-if="authState.isAdmin()" class="position-absolute top-0 end-0 m-3 d-flex gap-2">
        <RouterLink :to="`/admin/areas/${area.id}`" class="btn btn-primary shadow-sm">
          <i class="bi bi-pencil-square"></i>
        </RouterLink>

        <button class="btn btn-primary shadow-sm" v-on:click="deleteArea">
          <i class="bi bi-trash"></i>
        </button>
      </div>


      <div class="position-absolute bottom-0 start-0 p-4 text-white">
        <h1 class="mb-1">{{ area.name }}</h1>
        <p class="mb-0">
          {{ area.description }}
        </p>
      </div>
    </div>

    <div class="row g-4">
      <div class="col-lg-4">
        <div class="card h-100">
          <div class="card-body d-flex flex-column">
            <div class="d-flex align-items-center justify-content-between mb-3">
              <h5 class="card-title mb-0">Your Status</h5>

              <span class="badge" :class="completedAllTrainings ? 'text-bg-success' : 'text-bg-danger'">
                {{ completedAllTrainings ? 'Certified' : 'Incomplete' }}
              </span>
            </div>

            <p class="mb-1 fw-semibold">
              {{ completedAllTrainings ? 'All trainings complete' : 'Trainings required' }}
            </p>

            <small v-if="completedAllTrainings" class="text-body-secondary mb-3">
              You have completed all required trainings
            </small>
            <small v-else class="text-body-secondary mb-3">
              Complete all required trainings to gain access to {{ area.name }}
            </small>

            <div class="mt-auto">
              <div class="progress" style="height: 6px">
                <div class="progress-bar progress-bar-animated progress-bar-striped" role="progressbar"
                  aria-valuenow="0" aria-valuemin="0" aria-valuemax="100" :style="{ width: progress + '%' }"></div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="col-lg-8">
        <div class="card h-100">
          <div class="card-body">
            <h5 class="card-title mb-3">Required Trainings</h5>

            <ul class="list-group list-group-flush">
              <li v-if="trainings?.length === 0" class="list-group-item d-flex justify-content-between text-muted">
                No trainings associated with area
              </li>
              <li v-for="training in trainings" :key="training.id"
                class="list-group-item d-flex justify-content-between">
                <RouterLink :to="`/training/${training.id}`" class="training-link">
                  {{ training.title }}
                </RouterLink>

                <span v-if="userTrainings?.some((t) => t.trainingId === training.id)" class="badge text-bg-success">
                  Completed
                </span>
                <span v-else class="badge text-bg-warning"> Required </span>
              </li>
            </ul>
          </div>
        </div>
      </div>
    </div>

    <div class="card mt-4" v-if="resources && resources.length > 0">
      <div class="card-body">
        <h5 class="card-title">Resources</h5>
        <div class="d-flex gap-2 flex-wrap">
          <a v-for="r in resources" :href="r.url" :key="r.id" target="_blank">
            <button class="btn btn-outline-primary btn-sm">{{ r.name }}</button>
          </a>
        </div>
      </div>
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
.image-gradient {
  background: linear-gradient(to top,
      rgba(0, 0, 0, 0.8) 0%,
      rgba(0, 0, 0, 0.4) 40%,
      rgba(0, 0, 0, 0.1) 70%,
      rgba(0, 0, 0, 0) 100%);
}

.training-link {
  text-decoration: none;
  color: var(--bs-dark);
}

.training-link:hover {
  color: var(--bs-primary);
  text-decoration: underline;
}
</style>
