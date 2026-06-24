<template>
  <main class="container d-flex justify-content-center align-items-center">
    <div class="bg-text z-n1">404</div>

    <div class="content text-center z-1">
      <h1 class="fw-bold">Nothing here</h1>
      <p class="text-break text-muted">
        Nothing here go away
      </p>

      <RouterLink to="/" class="btn btn-primary">
        Return Home
      </RouterLink>
    </div>
  </main>
</template>

<style scoped>

main {
    margin-top: 5rem;
}

.bg-text {
  position: absolute;
  font-size: 18rem;
  font-weight: 900;
  opacity: 0.05;
  user-select: none;
}
</style>