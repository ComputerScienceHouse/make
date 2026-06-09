
<script setup lang="ts" generic="T extends Record<string, any>">
import { authState } from '@/auth';
import type { Area } from '@/models/areas';
import { computed} from 'vue';



interface Props<T> {
  data: T[];
}

const props = defineProps<Props<T>>();

const data = props.data

const cols = computed(() => {
    if (data?.length === 0 || !data) return [];
    return Object.keys(data[0]!) as (keyof Area)[]
})

console.log(data)

</script>
<template>
     <div class="table-responsive">
      <table class="table table-striped table-hover align-middle">
        <thead class="table-dark">
          <tr>
            <th v-for="col in cols" :key="col" class="text-nowrap">
              {{ col }}
            </th>
          </tr>
        </thead>

        <tbody>
          <tr v-for="(thing, i) in data" :key="i">
            <td
              v-for="col in cols"
              :key="col"
              class="text-truncate"
              style="max-width: 200px;"
            >
              {{ thing[col] }}
            </td>
          </tr>
        </tbody>
      </table>
    </div>
</template>