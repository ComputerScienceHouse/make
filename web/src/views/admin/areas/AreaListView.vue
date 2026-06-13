<script setup lang="ts">
import DynamicTable from "@/components/DynamicTable.vue";
import type { Area } from "@/models/areas";
import type { TableOptions } from '@/components/DynamicTable.vue'
import { ref, onMounted } from "vue";

const areas = ref<Area[]>([]);
const loading = ref(true);
const notFound = ref(false);

const tableOptions: TableOptions<Area> = {
    actions:  {
        edit: {
            // relative path, will redirect to
            // /admin/areas/id
            path: (area: Area) => `/admin/areas/${area.id}`
        }
    }
}

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
        <div class="d-flex justify-content-between align-items-center">
            <h1>Areas:</h1>
            <RouterLink to="/admin/areas/create" type="button" class="btn btn-primary">Create</RouterLink>
        </div>

        <DynamicTable :data="areas" :options="tableOptions"></DynamicTable>

    </main>
</template>

<style>
main {
    padding: 2rem;
    max-width: 100ch;
}
</style>
