<script setup lang="ts">
const sections = [
    {
        title: "Areas",
        description: "Manage access-controlled areas and their training requirements.",
        path: "/admin/areas",
        icon: "bi-door-open"
    },
    {
        title: "Create Area",
        description: "Create a new area.",
        path: "/admin/areas/create",
        icon: "bi-plus-circle"
    },
        {
        title: "Create Trainings",
        description: "Create a new training.",
        path: "/admin/trainings/create",
        icon: "bi-plus-circle"
    },
            {
        title: "Manage User Trainings ",
        description: "Manage trainings completed by a user",
        path: "/admin/trainings/create",
        icon: "bi-plus-circle"
    },
]
</script>

<template>
    <main class="container">
        <div class="container py-4">
            <h1 class="mb-4">Admin Dashboard</h1>

            <div class="row g-4">
                <div v-for="section in sections" :key="section.path" class="col-md-6 col-lg-4">
                    <RouterLink :to="section.path" class="text-decoration-none">
                        <div class="card h-100 shadow-sm">
                            <div class="card-body">
                                <h3 class="card-title">
                                    <i :class="`bi ${section.icon} me-2`"></i>
                                    {{ section.title }}
                                </h3>

                                <p class="card-text text-muted">
                                    {{ section.description }}
                                </p>
                            </div>
                        </div>
                    </RouterLink>
                </div>
            </div>
        </div>
    </main>

</template>