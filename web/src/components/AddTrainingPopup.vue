<script setup lang="ts">
import type { Training } from '@/models/trainings';
import { onMounted, ref } from 'vue';
import DynamicTable from './DynamicTable.vue';
import type { TableOptions } from './DynamicTable.vue';

const emit = defineEmits<{
    close: [];
}>();

function close() {
    emit('close');
}

interface Props {
    areaid: number;
}
const props = defineProps<Props>();


const tableOptions: TableOptions<Training> = {
    fields: {
        questions: { hidden: true },
    },
    actions: {
        checkbox: true,
    }
}

const loading = ref(true);
const trainings = ref<Training[]>();
const error = ref<string>();

// From table
const selectedTrainings = ref<Training[]>([]);


onMounted(async () => {
    try {
        loading.value = true

        const trainingRes = await fetch(`/api/trainings/`, {
            credentials: "include"
        })

        if (!trainingRes.ok) {
            throw new Error("Failed to fetch data");
        }

        trainings.value = await trainingRes.json()
    } catch (err) {
        console.error(err);
    } finally {
        loading.value = false;
    }
});

async function save() {
    const res = await fetch(`/api/trainings/area/${props.areaid}`, {
        method: "POST",
        body: JSON.stringify(selectedTrainings.value.map(s => s.id)),
        credentials: "include"
    })

    if (!res.ok) {
        error.value = `Failed to save area: ${res.status}`;
        throw new Error(`Failed to save area: ${res.status}`);
    } else {
        location.reload()
    }

    

    return;
}

</script>
<template>
    <div class="grayout"></div>
    <div class="modal d-block" tabindex="100">
        <div class="modal-dialog modal-dialog-centered">
            <div class="modal-content">
                <div class="modal-header">
                    <h5 class="modal-title">Add Training</h5>
                    <button type="button" class="btn-close" aria-label="Close" v-on:click="close"></button>
                </div>
                <div class="modal-body">
                    <p>Select trainings you wish to add</p>
                    <div v-if="loading">
                        Loading....
                    </div>
                    <div v-else-if="trainings">
                        <DynamicTable :data="trainings" :options="tableOptions" v-model="selectedTrainings"></DynamicTable>
                    </div>
                    <p v-if="error" class="text-danger">{{ error }}</p>
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-secondary" v-on:click="close">Close</button>
                    <button type="button" class="btn btn-primary" @click="save()">Save changes</button>
                </div>
            </div>
        </div>
    </div>

</template>

<style scoped>
.grayout {
    position: fixed;
    left: 0px;
    top: 0px;
    height: 100%;
    width: 100%;
    background-color: black;
    opacity: 0.5;
    z-index: 100;
}
</style>