<script setup lang="ts">
import type { Training, UserTraining } from '@/models/trainings'
import { onMounted, ref } from 'vue'

const emit = defineEmits<{
  close: []
}>()

function close() {
  emit('close')
}

interface Props {
  uuid: string
}
const props = defineProps<Props>()

const loading = ref(true)
const trainings = ref<Training[]>([])
const error = ref<string>()

// single selection
const selectedTrainingId = ref<number | null>(null)

// expiry date (YYYY-MM-DD for input[type="date"])
const expiresAt = ref<string>(
  new Date(Date.now() + 7 * 24 * 60 * 60 * 1000).toISOString().substring(0, 10),
)

onMounted(async () => {
  try {
    loading.value = true

    const trainingRes = await fetch(`/api/trainings/`, {
      credentials: 'include',
    })

    if (!trainingRes.ok) {
      throw new Error('Failed to fetch data')
    }

    trainings.value = await trainingRes.json()
  } catch (err) {
    console.error(err)
    error.value = 'Failed to load trainings'
  } finally {
    loading.value = false
  }
})

async function save() {
  if (!selectedTrainingId.value) return

  const payload: UserTraining = {
    userUuid: props.uuid,
    trainingId: selectedTrainingId.value,
    completedAt: new Date(),
    expiresAt: new Date(expiresAt.value),
  }

  const res = await fetch(`/api/user/${props.uuid}/trainings`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(payload),
    credentials: 'include',
  })

  if (!res.ok) {
    error.value = `Failed to save: ${res.status}`
    throw new Error(`Failed to save: ${res.status}`)
  }

  location.reload()
}
</script>

<template>
  <div class="grayout"></div>

  <div class="modal d-block" tabindex="100">
    <div class="modal-dialog modal-dialog-centered">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">Add Training</h5>
          <button type="button" class="btn-close" @click="close"></button>
        </div>

        <div class="modal-body">
          <p>Select a training and set expiry date</p>

          <div v-if="loading">Loading...</div>

          <div v-else>
            <!-- Dropdown -->
            <label class="form-label">Training</label>
            <select class="form-select" v-model="selectedTrainingId">
              <option :value="null" disabled>-- Select a training --</option>

              <option v-for="t in trainings" :key="t.id" :value="t.id">
                {{ t.title }} (id: {{ t.id }})
              </option>
            </select>

            <!-- Expiry -->
            <div class="mt-3" v-if="selectedTrainingId">
              <label class="form-label">Expiry date</label>
              <input type="date" class="form-control" v-model="expiresAt" />
            </div>

            <p v-if="error" class="text-danger mt-2">
              {{ error }}
            </p>
          </div>
        </div>

        <div class="modal-footer">
          <button class="btn btn-secondary" @click="close">Close</button>
          <button class="btn btn-primary" @click="save">Save changes</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.grayout {
  position: fixed;
  left: 0;
  top: 0;
  height: 100%;
  width: 100%;
  background-color: black;
  opacity: 0.5;
  z-index: 100;
}
</style>
