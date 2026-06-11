<script setup lang="ts" generic="T extends Record<string, any>">
import type { Area } from '@/models/areas';
import { computed } from 'vue';

export interface TableOptions<T> {
  actions?: {
    edit?: {
      path: (row: T) => string
    },
  }
}

interface Props<T> {
  data: T[];
  options: TableOptions<T>;
}

const props = defineProps<Props<T>>();

const data = props.data

const cols = computed(() => {
  if (data?.length === 0 || !data) return [];
  return Object.keys(data[0]!) as (keyof T)[]
})

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
        </tr>
      </thead>

      <tbody>
        <tr v-for="(row, i) in data" :key="i">
          <td v-for="col in cols" :key="col" class="text-truncate" style="max-width: 200px;">
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
        </tr>

      </tbody>
    </table>
  </div>
</template>