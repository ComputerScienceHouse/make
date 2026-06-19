<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useRoute } from "vue-router";
import DynamicForm from "@/components/DynamicForm.vue";
import type { FormOptions } from "@/components/DynamicForm.vue";
import router from "@/router";
import type { Training } from "@/models/trainings";
const training = ref<Training>();
const loading = ref(true);
const notFound = ref(false);

const route = useRoute();

const formOptions: FormOptions<Training> = {
    fields: {
        id: { 
            hidden: true 
        },
        questions: {
            type: "textarea"
        }
    }
}

onMounted(async () => {
    try {
        loading.value = true

        const trainingRes = await fetch(`/api/trainings/${route.params.id}`)

        if (trainingRes.status === 404) {
            notFound.value = true;
        }

        if (!trainingRes.ok) {
            throw new Error("Failed to fetch data");
        }

        training.value = await trainingRes.json()
    } catch (err) {
        console.error(err);
    } finally {
        loading.value = false;
    }
});

async function saveTraining(training: Training) {
    // TODO: error handling
    training.questions = JSON.parse(training.questions)

    const res = await fetch(`/api/trainings/${training.id}`, {
        method: "PUT",
        body: JSON.stringify(training),
        credentials: "include"
    })

    if (!res.ok) {
        throw new Error(`Failed to save area: ${res.status}`);
    }

    router.push({ path: `/admin/trainings` })

    return;
}
</script>


<template>
    <main class="container py-4" v-if="loading">
        <h1>Loading...</h1>
    </main>

    <main class="container" v-else-if="training">
        <div class="d-flex justify-content-between align-items-center">
            <h1>Editing Training "{{ training?.title }}"</h1>
            <span>ID: {{ training.id }}</span>
        </div>

        <DynamicForm :data="training" :options="formOptions" @submit="saveTraining" class="mb-3"></DynamicForm>
    </main>
</template>

<style scoped>
main {
    padding: 2rem;
    max-width: 100ch;
}
</style>
