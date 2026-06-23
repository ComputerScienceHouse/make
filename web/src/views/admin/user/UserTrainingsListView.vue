<script setup lang="ts">
import DynamicTable from '@/components/DynamicTable.vue'
import AddTrainingToUserPopup from '@/components/AddTrainingToUserPopup.vue'

import type { TableOptions } from '@/components/DynamicTable.vue'
import type { UserTraining } from '@/models/trainings'
import type { Member } from '@/models/member'

import { ref, onMounted } from 'vue'

const members = ref<Member[]>([])
const selectedMember = ref<Member>()
const memberSearch = ref('')

const trainings = ref<UserTraining[]>([])
const loading = ref(true)
const trainingsLoading = ref(false)

const showTrainingModal = ref(false)

const tableOptions: TableOptions<UserTraining> = {
  actions: {
    delete: {
      handler: async (training: UserTraining) => {
        if (!selectedMember.value) return

        const res = await fetch(
          `/api/user/${selectedMember.value.uuid}/trainings/${training.trainingId}`,
          {
            method: 'DELETE',
            credentials: 'include',
          },
        )

        if (!res.ok) {
          throw new Error('Failed to delete training')
        }
      },
    },
  },
  fields: {
    userUuid: {
      hidden: true,
    },
  },
}

async function loadTrainings() {
  const member = members.value.find((m) => m.username === memberSearch.value)

  if (!member) {
    alert('Please select a valid member.')
    return
  }

  selectedMember.value = member

  try {
    trainingsLoading.value = true

    const trainingsRes = await fetch(`/api/user/${member.uuid}/trainings`, {
      credentials: 'include',
    })

    if (!trainingsRes.ok) {
      throw new Error('Failed to fetch trainings')
    }

    trainings.value = await trainingsRes.json()
  } catch (err) {
    console.error(err)
  } finally {
    trainingsLoading.value = false
  }
}

onMounted(async () => {
  try {
    loading.value = true

    const membersRes = await fetch('/api/members/active', {
      credentials: 'include',
    })

    if (!membersRes.ok) {
      throw new Error('Failed to fetch members')
    }

    members.value = await membersRes.json()
  } catch (err) {
    console.error(err)
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <AddTrainingToUserPopup
    v-if="showTrainingModal && selectedMember"
    :uuid="selectedMember.uuid"
    @close="showTrainingModal = false"
  />

  <main class="container py-4" v-if="loading">
    <h1>Loading...</h1>
  </main>

  <main class="container" v-else>
    <h1 class="mb-4">Editing User Trainings</h1>
    <div class="mb-4">
      <label for="memberSearch" class="form-label"> Select member </label>

      <div class="input-group">
        <input
          id="memberSearch"
          class="form-control"
          list="memberOptions"
          v-model="memberSearch"
          placeholder="Type a username..."
        />

        <button class="btn btn-primary" @click="loadTrainings">Load Trainings</button>
      </div>

      <datalist id="memberOptions">
        <option v-for="member in members" :key="member.uuid" :value="member.username" />
      </datalist>
    </div>

    <div v-if="selectedMember">
      <div class="d-flex justify-content-between align-items-center mb-3">
        <h1>{{ selectedMember.username }}'s Completed Trainings</h1>

        <button class="btn btn-primary" @click="showTrainingModal = true">
          <i class="bi bi-plus-lg"></i>
        </button>
      </div>

      <div v-if="trainingsLoading">Loading trainings...</div>

      <DynamicTable v-else :data="trainings" :options="tableOptions" />
    </div>
  </main>
</template>

<style scoped>
main {
  padding: 2rem;
  max-width: 100ch;
}
</style>
