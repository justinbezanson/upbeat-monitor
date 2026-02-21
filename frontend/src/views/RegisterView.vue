<template>
  <div class="max-w-md w-full bg-white p-8 rounded-lg shadow-md">
    <h2 class="text-2xl font-bold text-center text-gray-800 mb-6">Register</h2>
    <form @submit.prevent="handleRegister">
      <div class="mb-4">
        <label for="email" class="block text-gray-700 text-sm font-bold mb-2">Email:</label>
        <Input
          type="email"
          id="email"
          v-model="email"
          required
        />
      </div>
      <div class="mb-6">
        <label for="password" class="block text-gray-700 text-sm font-bold mb-2">Password:</label>
        <Input
          type="password"
          id="password"
          v-model="password"
          required
        />
      </div>
      <div v-if="authStore.error" class="text-red-500 mb-4 text-center">
        {{ authStore.error }}
      </div>
      <div v-if="authStore.successMessage" class="text-green-500 mb-4 text-center">
        {{ authStore.successMessage }}
      </div>
      <div class="flex items-center justify-between">
        <Button
          type="submit"
          class="w-full"
          :disabled="authStore.isLoading"
        >
          {{ authStore.isLoading ? 'Registering...' : 'Register' }}
        </Button>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';
import { useRouter } from 'vue-router';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { useAuthStore } from '@/stores/auth'; // Import the auth store

const email = ref('');
const password = ref('');
const router = useRouter();
const authStore = useAuthStore(); // Use the auth store

const handleRegister = async () => {
  const success = await authStore.register(email.value, password.value);
  if (success) {
    // Optionally redirect to login or dashboard
    router.push('/login'); // Assuming a login page exists or will be created
  }
};

// Watch for changes in successMessage to potentially clear form or redirect
watch(() => authStore.successMessage, (newValue) => {
  if (newValue) {
    email.value = '';
    password.value = '';
  }
});

// Watch for changes in error to clear success message if an error occurs later
watch(() => authStore.error, (newValue) => {
  if (newValue) {
    authStore.successMessage = null; // Clear success message if an error appears
  }
});
</script>

<style scoped>
/* You can add component-specific styles here if needed */
</style>