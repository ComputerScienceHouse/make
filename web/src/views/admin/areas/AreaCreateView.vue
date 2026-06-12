<script setup lang="ts">
import type { Area } from "@/models/areas";
import { ref } from "vue";
import DynamicForm, { type FormOptions } from "@/components/DynamicForm.vue";
import router from "@/router";

const area = ref<Area>({
  id: 0, // id is ignored by api when creating a new area
  name: "",
  description: "",
  ldapGroup: "",
  photourl: ""
})

const formOptions: FormOptions<Area> = {
    fields: {
        id: { hidden: true },
    },
}

async function saveArea(area: Area) {
    const res = await fetch(`/api/areas/create`, {
        method: "POST",
        body: JSON.stringify(area),
        credentials: "include"
    })

    if (!res.ok) {
        throw new Error(`Failed to save area: ${res.status}`);
    }

    router.push({ path: `/areas/${area.id}` })

    return;
}
</script>


<template>
    <main class="container" v-if="area">
        <div class="d-flex justify-content-between align-items-center">
            <h1>Creating New Area</h1>
        </div>

        <DynamicForm :data="area" :options="formOptions" @submit="saveArea" class="mb-3"></DynamicForm>

    </main>
</template>

<style scoped>
main {
    padding: 2rem;
    max-width: 100ch;
}
</style>
