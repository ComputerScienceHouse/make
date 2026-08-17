<script setup lang="ts">
import { onMounted, ref } from 'vue'
import DynamicTable from './DynamicTable.vue'
import type { TableOptions } from './DynamicTable.vue'
import type { Resource } from '@/models/resources.ts'

const emit = defineEmits<{
  close: []
}>()

function close() {
  emit('close')
}

interface Props {
  areaid: number
}
const props = defineProps<Props>()

const tableOptions: TableOptions<Resource> = {
  actions: {
    checkbox: true,
  },
}

const loading = ref(true)
const resources = ref<Resource[]>()
const error = ref<string>()

// From table
const selectedResources = ref<Resource[]>([])

onMounted(async () => {
  try {
    loading.value = true

    const resourceRes = await fetch(`/api/resources/`, {
      credentials: 'include',
    })

    if (!resourceRes.ok) {
      throw new Error('Failed to fetch data')
    }

    resources.value = await resourceRes.json()
  } catch (err) {
    console.error(err)
    error.value = `Error while fetching resources`
  } finally {
    loading.value = false
  }
})

async function save() {
  const res = await fetch(`/api/resources/area/${props.areaid}`, {
    method: 'POST',
    body: JSON.stringify(selectedResources.value.map((s) => s.id)),
    credentials: 'include',
  })

  if (!res.ok) {
    error.value = `Failed to save area: ${res.status}`
    throw new Error(`Failed to save area: ${res.status}`)
  } else {
    location.reload()
  }

  return
}
</script>
<template>
  <div class="grayout"></div>
  <div class="modal d-block" tabindex="100">
    <div class="modal-dialog modal-dialog-centered">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">Add Resource</h5>
          <button type="button" class="btn-close" aria-label="Close" v-on:click="close"></button>
        </div>
        <div class="modal-body">
          <p>Select resources you wish to add</p>
          <div v-if="loading">Loading....</div>
          <div v-else-if="resources">
            <DynamicTable
              :data="resources"
              :options="tableOptions"
              v-model="selectedResources"
            ></DynamicTable>
          </div>
          <p v-if="error" class="text-danger">{{ error }}</p>
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-primary" @click="save()">Save changes</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.grayout {
  position: fixed;
  left: 0px;
  top: 0px;
  height: 100%;
  width: 100%;
  background-color: black;
  opacity: 0.5;
  z-index: 100;
}
</style>
