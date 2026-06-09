<script setup lang="ts">
import type { Area } from '@/models/areas';
import { authState } from "@/auth";

interface Props {
    areas: Area[]
}

const user = authState.user

const props = defineProps<Props>()
</script>

<template>
    <div class="card mb-4 shadow-sm">
      <div class="d-flex flex-row justify-content-between">
        <div class="card-body d-flex align-items-center gap-3">
          <img :src="'https://profiles.csh.rit.edu/image/' + user?.preferred_username" class="rounded-circle"
            height="100" width="100">
          <div>
            <h2 class="mb-0">{{ user?.name }}</h2>
            <span class="text-body-secondary">{{ user?.preferred_username }}</span>
          </div>
        </div>

        <div class="d-flex align-items-center pe-4">
          <div>
            <h6 class="mb-2 text-body-secondary">Access Status:</h6>
            <ul class="list-unstyled mb-0">
              <li v-for="area in props.areas" :key="area.id" class="d-flex align-items-center gap-2">

                <i :class="false ? 'bi bi-check-circle-fill text-success' : 'bi bi-x-circle-fill text-danger'"></i>

                {{ area.name }}
              </li>
            </ul>
          </div>
        </div>
      </div>
    </div>
</template>