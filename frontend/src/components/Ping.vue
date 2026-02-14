<script setup lang="ts">
import { ref, onMounted } from 'vue'

const pingResult = ref<string | null>(null)
const loading = ref(true)
const error = ref<string | null>(null)

const fetchPing = async () => {
  try {
    loading.value = true
    error.value = null
    // Use the VITE_API_BASE_URL environment variable
    const response = await fetch(`${import.meta.env.VITE_API_BASE_URL}/ping`)
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }
    const data = await response.text() // Assuming ping returns plain text
    pingResult.value = data
  } catch (e: any) {
    error.value = `Failed to fetch ping: ${e.message}`
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchPing()
})
</script>

<template>
  <div class="border border-gray-300 p-4 mt-5 rounded-lg bg-gray-50">
    <h2 class="text-red-800 text-lg font-semibold mb-2">Ping API Status</h2>
    <p v-if="loading">Loading...</p>
    <p v-else-if="error" class="text-red-600 font-bold">{{ error }}</p>
    <p v-else-if="pingResult">Result: <span class="font-bold text-blue-600">{{ pingResult }}</span></p>
    <p v-else>No ping result available.</p>
    <button
      @click="fetchPing"
      :disabled="loading"
      class="bg-blue-600 text-white py-2 px-4 rounded-md cursor-pointer transition-colors duration-300
             hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed"
    >
      Refresh Ping
    </button>
  </div>
</template>
