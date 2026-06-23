<script setup lang="ts">
const sections = [
    {
        title: "Areas",
        items: [
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
            }
        ]

    },
    {
        title: "Trainings",
        items: [
            {
                title: "Trainings",
                description: "Manage trainings.",
                path: "/admin/trainings",
                icon: "bi-rulers"
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
                path: "/admin/user/trainings",
                icon: "bi-plus-circle"
            },
        ]

    },
]
</script>

<template>
    <main class="container">
        <div class="container py-4">
            <h1 class="mb-4">Admin Dashboard</h1>

            <div>
                <div v-for="section in sections" :key="section.title" class="row g-4 mb-2">
                    <h1> {{ section.title }}</h1>
                    <div v-for="item in section.items" :key="item.title" class="col-md-6 col-lg-4 mt-2 mb-2">
                        <RouterLink :to="item.path" class="text-decoration-none">
                            <div class="card h-100 shadow-sm">
                                <div class="card-body">
                                    <h3 class="card-title">
                                        <i :class="`bi ${item.icon} me-2`"></i>
                                        {{ item.title }}
                                    </h3>

                                    <p class="card-text text-muted">
                                        {{ item.description }}
                                    </p>
                                </div>
                            </div>
                        </RouterLink>
                    </div>
                </div>
            </div>
        </div>
    </main>

</template>