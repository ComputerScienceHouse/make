<script setup lang="ts">
import DynamicTable from "@/components/DynamicTable.vue";
import type { Area } from "@/models/areas";
import { ref, onMounted } from "vue";

const areas = ref<Area[]>([]);
const loading = ref(true);
const notFound = ref(false);

onMounted(async () => {
    try {
        loading.value = true

        const areaRes = await fetch(`/api/areas/`);

        if (areaRes.status === 404) {
            notFound.value = true;
        }

        if (!areaRes.ok) {
            throw new Error("Failed to fetch data");
        }

        areas.value = await areaRes.json()
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

    <main class="container" v-else>
        <h1>Areas:</h1>
        <DynamicTable :data="areas"></DynamicTable>
    </main>
</template>

<style>
main {
    padding: 2rem;
    max-width: 100ch;
}
</style>
