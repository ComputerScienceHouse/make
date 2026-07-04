<script setup lang="ts">
import { RouterView } from 'vue-router'
import AppNavbar from '@/components/AppNavbar.vue'
import AppFooter from './components/AppFooter.vue'
</script>

<template>
  <AppNavbar />

  <RouterView />

  <AppFooter></AppFooter>
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
