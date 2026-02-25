<template>
  <div class="flex min-h-screen flex-1 flex-col justify-center px-6 py-12 lg:px-8 bg-gray-900 w-full m-0">
    <div class="sm:mx-auto sm:w-full sm:max-w-sm">
      <img class="mx-auto h-10 w-auto" src="https://tailwindcss.com/plus-assets/img/logos/mark.svg?color=indigo&shade=500" alt="Your Company" />
      <h2 class="mt-10 text-center text-2xl/9 font-bold tracking-tight text-white">Set new password</h2>
      <p class="mt-2 text-center text-sm text-gray-400">
        Please enter your new password below.
      </p>
    </div>

    <div class="mt-10 sm:mx-auto sm:w-full sm:max-w-sm">
      <form class="space-y-6" @submit.prevent="handleResetPassword">
        <div>
          <label for="password" class="block text-sm/6 font-medium text-gray-100">New Password</label>
          <div class="mt-2">
            <Input 
              type="password" 
              name="password" 
              id="password" 
              autocomplete="new-password" 
              required 
              v-model="password"
              class="block w-full rounded-md bg-white/5 px-3 py-1.5 text-base text-white outline-1 -outline-offset-1 outline-white/10 placeholder:text-gray-500 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-500 sm:text-sm/6" 
            />
          </div>
        </div>

        <div>
          <label for="confirmPassword" class="block text-sm/6 font-medium text-gray-100">Confirm New Password</label>
          <div class="mt-2">
            <Input 
              type="password" 
              name="confirmPassword" 
              id="confirmPassword" 
              autocomplete="new-password" 
              required 
              v-model="confirmPassword"
              class="block w-full rounded-md bg-white/5 px-3 py-1.5 text-base text-white outline-1 -outline-offset-1 outline-white/10 placeholder:text-gray-500 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-500 sm:text-sm/6" 
            />
          </div>
        </div>

        <div v-if="authStore.error" class="text-red-500 text-sm text-center">
          {{ authStore.error }}
        </div>

        <div v-if="authStore.successMessage" class="text-green-500 text-sm text-center">
          {{ authStore.successMessage }}
        </div>

        <div>
          <Button 
            type="submit" 
            :disabled="authStore.isLoading || !token"
            class="flex w-full justify-center rounded-md bg-indigo-500 px-3 py-1.5 text-sm/6 font-semibold text-white hover:bg-indigo-400 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-500"
          >
            {{ authStore.isLoading ? 'Resetting...' : 'Reset password' }}
          </Button>
        </div>
      </form>

      <div v-if="!token" class="mt-4 text-center text-red-400 text-sm">
        Missing reset token. Please check your email link.
      </div>

      <p class="mt-10 text-center text-sm/6 text-gray-400">
        Back to
        {{ ' ' }}
        <router-link to="/login" class="font-semibold text-indigo-400 hover:text-indigo-300">Sign in</router-link>
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useAuthStore } from '@/stores/auth';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';

const route = useRoute();
const router = useRouter();
const authStore = useAuthStore();

const token = ref('');
const password = ref('');
const confirmPassword = ref('');

onMounted(() => {
  const t = route.query.token;
  if (t && typeof t === 'string') {
    token.value = t;
  }
});

const handleResetPassword = async () => {
  if (password.value !== confirmPassword.value) {
    authStore.error = "Passwords do not match.";
    return;
  }

  const success = await authStore.resetPassword(token.value, password.value);
  if (success) {
    // Optionally redirect to login after a delay
    setTimeout(() => {
        router.push('/login');
    }, 3000);
  }
};
</script>

<style scoped>
</style>
