<script setup lang="ts">
import type { Area } from "@/models/areas";
import { ref } from "vue";
import DynamicForm, { type FormOptions } from "@/components/DynamicForm.vue";
import router from "@/router";
import type { Training } from "@/models/trainings";

const training = ref<Training>({
  id: 0, // id is ignored by api when creating a new area
  title: "",
  description: "",
  questions: "",
})

const formOptions: FormOptions<Training> = {
    fields: {
        id: { hidden: true },
        questions: {
            type: "textarea"
        }
    },
}

async function saveArea(t: Training) {
    const res = await fetch(`/api/trainings/create`, {
        method: "POST",
        body: JSON.stringify(t),
        credentials: "include"
    })

    if (!res.ok) {
        throw new Error(`Failed to save training: ${res.status}`);
    }

    router.push({ path: `/admin/trainings/` })

    return;
}
</script>


<template>
    <main class="container" v-if="training">
        <div class="d-flex justify-content-between align-items-center">
            <h1>Creating New Training</h1>
        </div>

        <DynamicForm :data="training" :options="formOptions" @submit="saveArea" class="mb-3"></DynamicForm>

    </main>
</template>

<style scoped>
main {
    padding: 2rem;
    max-width: 100ch;
}
</style>
