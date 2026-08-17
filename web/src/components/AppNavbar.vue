<script setup lang="ts">
import { authState } from '@/auth'

const user = authState.user

function logout() {
  window.location.href = '/auth/logout'
}
</script>

<template>
  <nav class="navbar fixed-top navbar-expand-lg navbar-dark bg-primary">
    <div class="container">
      <RouterLink class="navbar-brand" to="/">
        <img
          class="object-fit-contain ms-2 me-2"
          height="32"
          src="https://assets.csh.rit.edu/pubsite/csh_logo_square.svg"
        />
        Make
      </RouterLink>
      <button
        class="navbar-toggler navbar-toggler-right"
        type="button"
        data-bs-toggle="collapse"
        data-bs-target="#navbarResponsive"
      >
        <span class="navbar-toggler-icon"></span>
      </button>
      <div class="collapse navbar-collapse" id="navbarResponsive">
        <ul class="navbar-nav mr-auto">
          <RouterLink v-if="authState.isAdmin()" class="nav-link" to="/admin">Admin</RouterLink>
        </ul>
        <ul v-if="user" class="nav navbar-nav ms-auto">
          <div class="nav-item navbar-user dropdown">
            <a
              class="nav-link dropdown-toggle"
              id="userDropdownLink"
              data-bs-toggle="dropdown"
              aria-expanded="false"
              role="button"
            >
              <img
                class="rounded-circle pfp me-2"
                :src="'https://profiles.csh.rit.edu/image/' + user.preferred_username"
              />
              <span class="me-1">{{ user.name }} </span>
              <span class="caret"></span>
            </a>
            <div class="dropdown-menu">
              <a class="dropdown-item" @click="logout()">Logout</a>
            </div>
          </div>
        </ul>
      </div>
    </div>
  </nav>
</template>

<style scoped>
.pfp {
  width: 2rem;
  height: 2rem;
}
</style>
