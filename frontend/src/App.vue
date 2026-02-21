<script setup lang="ts">
// No direct component imports here as they will be handled by routing
</script>

<template>
  <nav class="p-4 bg-gray-800 text-white flex justify-center space-x-4">
    <router-link to="/" class="hover:text-gray-300">Home</router-link>
    <router-link to="/register" class="hover:text-gray-300">Register</router-link>
  </nav>
  <div class="flex flex-col items-center justify-center min-h-screen p-8">
    <router-view />
  </div>
</template>

<style scoped>
/* Scoped styles if any */
</style>

