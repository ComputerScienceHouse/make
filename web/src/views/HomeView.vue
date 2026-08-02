<script setup lang="ts">
import { ref, onMounted } from 'vue'
import type { Area } from '@/models/areas'
import { authState } from '@/auth'

const user = authState.user
const areas = ref<Area[]>([])
const areasWithAccess = ref<number[]>([])

const error = ref('')
const loading = ref(true)

onMounted(async () => {
  try {
    loading.value = true

    const [areaRes, accessRes] = await Promise.all([
      fetch('/api/areas'),
      fetch(`/api/user/${user?.uuid}/areaAccess`),
    ])
    areas.value = await areaRes.json()
    areasWithAccess.value = await accessRes.json()
  } catch (err) {
    error.value = 'Error while fetching areas'
    console.error(err)
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <main class="container">
    <h1>Welcome to MAKE</h1>

    <p>The central hub for training and resources for CSH's special-use rooms.</p>
    <h1>Profile</h1>

    <div class="card mb-4 shadow-sm pe-none">
      <div class="d-flex flex-column flex-sm-row justify-content-between">
        <div class="card-body d-flex align-items-center gap-3">
          <img
            :src="'https://profiles.csh.rit.edu/image/' + user?.preferred_username"
            class="rounded-circle"
            height="100"
            width="100"
          />
          <div>
            <h2 class="mb-0">{{ user?.name }}</h2>
            <span class="text-body-secondary">{{ user?.preferred_username }}</span>
          </div>
        </div>

        <div class="d-flex align-items-center justify-content-center pb-3 pb-sm-0 text-center pe-md-4 mt-sm-2">
          <div>
            <h6 class="mb-2 text-body-secondary md:text-start">Access Status:</h6>
            <ul class="list-unstyled mb-0">
              <li v-for="area in areas" :key="area.id" class="d-flex align-items-center gap-2">
                <i
                  :class="
                    areasWithAccess.includes(area.id)
                      ? 'bi bi-check-circle-fill text-success'
                      : 'bi bi-x-circle-fill text-danger'
                  "
                ></i>

                {{ area.name }}
              </li>
            </ul>
          </div>
        </div>
      </div>
    </div>

    <div class="d-flex justify-content-between align-items-center">
      <h1>Areas:</h1>
      <RouterLink v-if="authState.isAdmin()" to="/admin/areas/create" class="btn btn-primary">
        <i class="bi-plus-lg"></i>
      </RouterLink>
    </div>

    <div class="row g-4">
      <div v-for="area in areas" :key="area.id" class="col-12 col-md-6">
        <RouterLink class="no_underline" :to="'/areas/' + area.id">
          <div class="card h-100 border-1 shadow-none">
            <img
              :src="area.photourl"
              class="card-img-top object-fit-cover"
              height="150"
              :alt="area.name"
            />

            <div class="card-body d-flex flex-column">
              <h4 class="card-title mb-2">
                {{ area.name }}
              </h4>

              <p class="card-text text-body-secondary flex-grow-1">
                {{ area.description }}
              </p>
            </div>
          </div>
        </RouterLink>
      </div>

      <div v-if="loading">
        <LoadingScreen></LoadingScreen>
      </div>

      <div v-else-if="error" class="text-danger">{{ error }}</div>

      <div v-else-if="areas.length === 0" class="text-muted">No areas found.</div>
    </div>
  </main>
</template>

<style scoped>
main {
  padding: 2rem;
  max-width: 100ch;
}

.card:hover {
  border-color: var(--bs-primary);
  cursor: pointer;
}

button {
  margin-top: 1rem;
  padding: 0.5rem 1rem;
  cursor: pointer;
}

.no_underline {
  text-decoration: none;
}
</style>
