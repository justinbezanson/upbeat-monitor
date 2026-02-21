<template>
  <div class="max-w-md w-full bg-white p-8 rounded-lg shadow-md">
    <h2 class="text-2xl font-bold text-center text-gray-800 mb-6">Login</h2>
    <form @submit.prevent="handleLogin">
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
      <div class="flex items-center justify-between">
        <Button
          type="submit"
          class="w-full"
          :disabled="authStore.isLoading"
        >
          {{ authStore.isLoading ? 'Logging in...' : 'Login' }}
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

console.log('LoginView component is being rendered.'); // Debug statement

const email = ref('');
const password = ref('');
const router = useRouter();
const authStore = useAuthStore(); // Use the auth store

const handleLogin = async () => {
  const success = await authStore.login(email.value, password.value);
  if (success) {
    // Redirect to a dashboard or home page after successful login
    router.push('/'); // Or '/dashboard' once it's created
  }
};

// Clear form fields on successful login (if desired)
watch(() => authStore.userToken, (newValue) => {
  if (newValue) {
    email.value = '';
    password.value = '';
  }
});

// Clear error message if successful login
watch(() => authStore.userToken, (newValue) => {
  if (newValue) {
    authStore.error = null;
  }
});
</script>

<style scoped>
/* You can add component-specific styles here if needed */
</style>
