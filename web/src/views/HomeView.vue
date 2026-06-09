<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useRouter } from "vue-router";
import type { Area } from "@/models/areas";
import QuickInfoCard from "@/components/QuickInfoCard.vue";



const areas = ref<Area[]>([]);
const router = useRouter();

onMounted(async () => {
  const response = await fetch("/api/areas");
  areas.value = await response.json();
});

function openArea(id: number) {
  router.push(`/areas/${id}`)
}
</script>

<template>
  <main class="container">
    <h1>Welcome to MAKE</h1>

    <p>
      The central hub for training and resources for CSH's special-use rooms.
    </p>
    <h1>Profile</h1>

    <QuickInfoCard :areas="areas"></QuickInfoCard>

    <h1>Areas</h1>

    <div class="row g-4">


      <div v-for="area in areas" :key="area.id" class="col-12 col-md-6">
        <RouterLink class="no_underline" :to="'/areas/' + area.id">
          <div class="card h-100 border-1 shadow-none">
            <img :src="area.photourl" class="card-img-top"
              height="150" :alt="area.name" />

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

      <div v-if="areas.length === 0">
        No areas found, add some!!nh nb ,mnbm,nbmbm,nbmbvnmbvn nmhjytrftrf
      </div>
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

.no_underline { text-decoration: none; }


</style>