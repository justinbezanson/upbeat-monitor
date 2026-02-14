<script setup lang="ts">
import HelloWorld from './components/HelloWorld.vue'
import Ping from './components/Ping.vue' // Import the Ping component
</script>

<template>
  <div class="flex flex-col items-center justify-center min-h-screen p-8">
    <div class="flex mb-8">
      <a href="https://vite.dev" target="_blank" class="block transition-all duration-300 hover:drop-shadow-[0_0_2em_#646cffaa]">
        <img src="/vite.svg" class="h-24 p-6" alt="Vite logo" />
      </a>
      <a href="https://vuejs.org/" target="_blank" class="block transition-all duration-300 hover:drop-shadow-[0_0_2em_#42b883aa]">
        <img src="./assets/vue.svg" class="h-24 p-6" alt="Vue logo" />
      </a>
    </div>
    <HelloWorld msg="Vite + Vue" />
    <Ping /> <!-- Add the Ping component here -->
  </div>
</template>
