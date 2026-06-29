<script setup lang="ts" generic="T extends Record<string, any>">
import { computed, ref } from 'vue'

export interface TableOptions<T> {
  actions?: {
    edit?: {
      path: (row: T) => string
    }
    delete?: {
      handler: (row: T) => Promise<void>
    }
    checkbox?: boolean
  }

  fields?: {
    [K in keyof T]?: {
      label?: string
      hidden?: boolean
    }
  }
}

interface Props<T> {
  data: T[]
  options: TableOptions<T>
}

const props = defineProps<Props<T>>()

const selectedRows = defineModel<T[]>({ default: [] })
function toggleRow(row: T, checked: boolean) {
  if (checked) {
    if (!selectedRows.value.includes(row)) {
      selectedRows.value.push(row)
    }
  } else {
    const index = selectedRows.value.indexOf(row)

    if (index !== -1) {
      selectedRows.value.splice(index, 1)
    }
  }
}

const data = props.data
const options = props.options

const error = ref<string>()

const cols = computed(() => {
  if (data?.length === 0 || !data) return []

  const keys = Object.keys(data[0]!) as (keyof T)[]
  return keys.filter((col) => !options.fields?.[col]?.hidden)
})

async function deleteRow(row: T) {
  try {
    await options.actions?.delete?.handler(row)

    const idx = data.indexOf(row)
    if (idx !== -1) {
      data.splice(idx, 1)
    }
  } catch (err) {
    error.value = 'Error while deleting row'
    console.error(err)
  }
}
</script>
<template>
  <div class="table-responsive">
    <table class="table table-striped table-hover align-middle">
      <thead class="table-dark">
        <tr>
          <th v-for="col in cols" :key="col" class="text-nowrap">
            {{ col }}
          </th>
          <th v-if="options.actions?.edit"></th>
          <th v-if="options.actions?.delete"></th>
          <th v-if="options.actions?.checkbox"></th>
        </tr>
      </thead>

      <tbody>
        <tr v-for="(row, i) in data" :key="i">
          <td v-for="col in cols" :key="col" class="text-truncate" style="max-width: 200px">
            {{ row[col] }}
          </td>

          <td v-if="options.actions?.edit">
            <div class="d-flex justify-content-center">
              <RouterLink :to="options.actions.edit.path(row)">
                <button type="button" class="btn btn-primary">
                  <i class="bi bi-pencil-square"></i>
                </button>
              </RouterLink>
            </div>
          </td>

          <td v-if="options.actions?.delete">
            <div class="d-flex justify-content-center">
              <button type="button" class="btn btn-primary" @click="deleteRow(row)">
                <i class="bi bi-trash"></i>
              </button>
            </div>
          </td>

          <td v-if="options.actions?.checkbox">
            <input
              type="checkbox"
              :checked="selectedRows.includes(row)"
              class="form-check-input"
              @change="toggleRow(row, ($event.target as HTMLInputElement).checked)"
            />
          </td>
        </tr>
      </tbody>
    </table>

    <span v-if="cols.length === 0"> No recorded data </span>
    <p v-if="error" class="text-danger">{{ error }}</p>
  </div>
</template>
