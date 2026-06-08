<script setup lang="ts">
import type { Area } from "@/models/areas";
import { ref, onMounted } from "vue";
import { useRoute } from "vue-router";

const area = ref<Area | null>(null);
const loading = ref(true);
const notFound = ref(false);

const route = useRoute();

onMounted(async () => {
    try {
        const response = await fetch(
            `/api/areas/${route.params.id}`
        );

        if (response.ok) {
            area.value = await response.json();
        } else if (response.status === 404) {
            notFound.value = true;
        } else {
            throw new Error("Failed to fetch area");
        }
    } catch (err) {
        console.error(err);
    } finally {
        loading.value = false;
    }
});
</script>

<template>
    <main class="container py-4" v-if="loading">
        <h1>Loading...</h1>
    </main>

    <main class="container py-4" v-else-if="area">
        <div class="position-relative rounded overflow-hidden mb-4" style="height: 280px;">
            <img :src="area.photourl" :alt="area.name"
                class="w-100 h-100 object-fit-cover" />

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
                    <div class="card-body">
                        <h5 class="card-title">Your Status</h5>

                        <p class="mb-2">
                            <span class="badge text-bg-success">
                                Certified
                            </span>
                        </p>

                        <small class="text-body-secondary">
                            You have completed all the bs
                        </small>
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
                            <li class="list-group-item d-flex justify-content-between">
                                general shit
                                <span class="badge text-bg-success">
                                    Complete
                                </span>
                            </li>

                            <li class="list-group-item d-flex justify-content-between">
                                lab shit
                                <span class="badge text-bg-success">
                                    Complete
                                </span>
                            </li>

                            <li class="list-group-item d-flex justify-content-between">
                                soldering shit
                                <span class="badge text-bg-warning">
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