<script setup lang="ts">
import { RouterView } from 'vue-router'
import AppNavbar from '@/components/AppNavbar.vue'
import AppFooter from './components/AppFooter.vue'

const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')

function updateTheme() {
  document.documentElement.setAttribute(
    'data-bs-theme',
    mediaQuery.matches ? 'dark' : 'light'
  )
}

updateTheme()
mediaQuery.addEventListener('change', updateTheme)
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
