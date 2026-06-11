<script setup lang="ts">
import DynamicTable from "@/components/DynamicTable.vue";
import type { Area } from "@/models/areas";
import type { TableOptions } from '@/components/DynamicTable.vue'
import { ref, onMounted } from "vue";
import { useRoute } from "vue-router";
import DynamicForm, { type FormOptions } from "@/components/DynamicForm.vue";
import router from "@/router";

const area = ref<Area>();
const loading = ref(true);
const notFound = ref(false);

const route = useRoute();

const tableOptions: FormOptions<Area> = {
    fields: {
        id: { hidden: true },
    }
}

onMounted(async () => {
    try {
        loading.value = true

        const areaRes = await fetch(`/api/areas/${route.params.id}`);

        if (areaRes.status === 404) {
            notFound.value = true;
        }

        if (!areaRes.ok) {
            throw new Error("Failed to fetch data");
        }

        area.value = await areaRes.json()
    } catch (err) {
        console.error(err);
    } finally {
        loading.value = false;
    }
});

async function saveArea(area: Area) {
    const res = await fetch(`/api/areas/${area.id}`, {
        method: "PUT",
        body: JSON.stringify(area),
        credentials: "include"
    })

    if (!res.ok) {
        throw new Error(`Failed to save area: ${res.status}`);
    }

    router.push({path: `/areas/${area.id}`})

    return;
}
</script>


<template>
    <main class="container py-4" v-if="loading">
        <h1>Loading...</h1>
    </main>

    <main class="container" v-else>
        <div class="d-flex justify-content-between align-items-center">
            <h1>Editing "{{ area?.name }}"</h1>
        </div>

        <DynamicForm :data="area!" :options="tableOptions" @submit="saveArea"></DynamicForm>

    </main>
</template>

<style>
main {
    padding: 2rem;
    max-width: 100ch;
}
</style>
