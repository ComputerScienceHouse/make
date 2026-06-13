<script setup lang="ts">
import DynamicTable from "@/components/DynamicTable.vue";
import type { Area } from "@/models/areas";
import type { TableOptions } from '@/components/DynamicTable.vue'
import { ref, onMounted, computed } from "vue";
import type { Training } from "@/models/trainings";
import { useRoute } from "vue-router";
import AddTrainingToUserPopup from "@/components/AddTrainingToUserPopup.vue";

const trainings = ref<Training[]>([]);
const loading = ref(true);
const route = useRoute();
const showTrainingModal = ref(false);

const uuid = computed(() => {
    return route.params.uuid as string
})

const tableOptions: TableOptions<Training> = {
    actions:  {
        delete: {
            handler: async (training: Training) => {
                await fetch(`/api/user/${uuid}/${training.id}`, {
                    method: "DELETE",
                    credentials: "include"
                })
            }
        }
    }
}

onMounted(async () => {
    try {
        loading.value = true
        const trainingsRes = await fetch(`/api/user/${route.params.uuid}/trainings`);

        if (!trainingsRes.ok) {
            throw new Error("Failed to fetch data");
        }

        trainings.value = await trainingsRes.json()
    } catch (err) {
        console.error(err);
    } finally {
        loading.value = false;
    }
});
</script>


<template>
    <AddTrainingToUserPopup v-if="showTrainingModal && trainings" :uuid="uuid"  @close="showTrainingModal = false">
    </AddTrainingToUserPopup>

    <main class="container py-4" v-if="loading">
        <h1>Loading...</h1>
    </main>

    <main class="container" v-else>
        <div class="d-flex justify-content-between align-items-center">
            <h1>User Trainings:</h1>
            <button to="/admin/trainings/create" class="btn btn-primary" @click="showTrainingModal = true">Add</button>
        </div>

        <DynamicTable :data="trainings" :options="tableOptions"></DynamicTable>

    </main>
</template>

<style>
main {
    padding: 2rem;
    max-width: 100ch;
}
</style>
