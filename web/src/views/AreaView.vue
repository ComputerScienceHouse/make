<script setup lang="ts">
import { authState } from "@/auth";
import type { Area } from "@/models/areas";
import type { Training } from "@/models/trainings";
import { ref, onMounted, computed } from "vue";
import { useRoute } from "vue-router";

const user = authState.user;

const area = ref<Area | null>(null);
const trainings = ref<Training[] | null>(null);
const userTrainings = ref<number[] | null>(null);
const loading = ref(true);
const notFound = ref(false);

const completedAllTrainings = computed(() => {
    return trainings.value?.every(t => userTrainings.value?.includes(t.id)) || false
})

const progress = computed(() => {
  if (!trainings.value?.length) return 0;
  return (userTrainings.value?.length / trainings.value.length) * 100;
});
const route = useRoute();

onMounted(async () => {
    try {

        const [areaRes, trainingsRes, userTrainingsRes] = await Promise.all([
            fetch(`/api/areas/${route.params.id}`),
            fetch(`/api/trainings/area/${route.params.id}`),
            fetch(`/api/user/${user?.uuid}/trainings`),
        ])

        if (areaRes.status === 404) {
            notFound.value = true;
        }

        if (!areaRes.ok || !trainingsRes.ok || !userTrainingsRes.ok) {
            throw new Error("Failed to fetch data");
        }

        area.value = await areaRes.json()
        trainings.value = await trainingsRes.json()
        userTrainings.value = await userTrainingsRes.json()
    } catch (err) {
        console.error(err);
    } finally {
        loading.value = false;
        console.log(userTrainings.value)
    }
});
</script>

<template>
    <main class="container py-4" v-if="loading">
        <h1>Loading...</h1>
    </main>

    <main class="container py-4" v-else-if="area">
        <div class="position-relative rounded overflow-hidden mb-4" style="height: 280px;">
            <img :src="area.photourl" :alt="area.name" class="w-100 h-100 object-fit-cover" />

            <div class="position-absolute top-0 start-0 w-100 h-100 image-gradient"></div>

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

                            <span class="badge rounded-pill px-3 py-2"
                                :class="completedAllTrainings ? 'text-bg-success' : 'text-bg-danger'">
                                {{ completedAllTrainings ? 'Certified' : 'Incomplete' }}
                            </span>
                        </div>

                        <p class="mb-1 fw-semibold">
                            {{ completedAllTrainings ? 'All training complete' : 'Training required' }}
                        </p>

                        <small v-if="completedAllTrainings" class="text-body-secondary mb-3">
                            You have completed all the bs
                        </small>
                        <small v-else class="text-body-secondary mb-3">
                            Complete all required training to gain access {{ area.name }}
                        </small>

                        <div class="mt-auto">
                            <div class="progress" style="height: 6px;">
                                <div class="progress-bar" role="progressbar" aria-valuenow="0" aria-valuemin="0"
                                    aria-valuemax="100" :style="{ width: progress + '%' }"></div>
                            </div>
                        </div>


                    </div>
                </div>
            </div>




            <div class="col-lg-8">
                <div class="card h-100">
                    <div class="card-body">
                        <h5 class="card-title mb-3">
                            Required Training
                        </h5>

                        <ul class="list-group list-group-flush">
                            <li v-for="training in trainings" :key="training.id" class="list-group-item d-flex justify-content-between">
                                {{ training.title }}
                                <span v-if="userTrainings?.includes(training.id)" class="badge text-bg-success">
                                    Completed
                                </span>
                                <span v-else class="badge text-bg-warning">
                                    Required
                                </span>
                            </li>
                        </ul>
                    </div>
                </div>
            </div>
        </div>

        <div class="card mt-4">
            <div class="card-body">
                <h5 class="card-title">Resources</h5>

                <div class="d-flex gap-2 flex-wrap">
                    <button class="btn btn-outline-primary btn-sm">
                        Safety Rules
                    </button>

                    <button class="btn btn-outline-primary btn-sm">
                        Equipment Guide
                    </button>
                </div>
            </div>
        </div>
    </main>

    <main class="container py-4" v-else-if="notFound">
        <h1>Area Not Found</h1>
        <p class="text-body-secondary">
            The requested workshop does not exist.
        </p>
    </main>

    <main class="container py-4" v-else>
        <h1>Something went wrong.</h1>
        <p class="text-body-secondary">
            Unable to load this area right now.
        </p>
    </main>
</template>


<style scoped>
.image-gradient {
    background:
        linear-gradient(to top,
            rgba(0, 0, 0, .8) 0%,
            rgba(0, 0, 0, .4) 40%,
            rgba(0, 0, 0, .1) 70%,
            rgba(0, 0, 0, 0) 100%);
}
</style>