<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
import Navbar from './components/Navbar.vue';
</script>

<template>
  <Navbar />

  <RouterView />
</template>

<style>
main {
  margin-top: 1em;
  margin-left: 2em;
  margin-right: 2em;
}

body {
  padding-top: 3rem;
}
</style>