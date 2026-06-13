<script setup lang="ts">
const version = git_describe
const year = new Date().getFullYear()
</script>

<template>
  <footer class="fixed-bottom bg-body text-body-secondary py-2">
    <div class="container d-flex justify-content-center align-items-center gap-3 small">
      <span>
        <i class="bi bi-c-circle me-1"></i>
        {{ year }} Computer Science House
      </span>

      <span class="text-muted">|</span>

      <span>
        <i class="bi bi-git me-1"></i>
        {{ version }}
      </span>

      <span class="text-muted">|</span>

      <a href="https://github.com/ComputerScienceHouse/make" target="_blank"
        class="text-decoration-none text-body-secondary">
        <i class="bi bi-github me-1"></i>
        Source
      </a>
    </div>
  </footer>
</template>