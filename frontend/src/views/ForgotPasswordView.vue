<template>
  <div class="flex min-h-screen flex-1 flex-col justify-center px-6 py-12 lg:px-8 bg-gray-900 w-full m-0">
    <div class="sm:mx-auto sm:w-full sm:max-w-sm">
      <img class="mx-auto h-10 w-auto" src="https://tailwindcss.com/plus-assets/img/logos/mark.svg?color=indigo&shade=500" alt="Your Company" />
      <h2 class="mt-10 text-center text-2xl/9 font-bold tracking-tight text-white">Reset your password</h2>
      <p class="mt-2 text-center text-sm text-gray-400">
        Enter your email address and we'll send you a link to reset your password.
      </p>
    </div>

    <div class="mt-10 sm:mx-auto sm:w-full sm:max-w-sm">
      <form class="space-y-6" @submit.prevent="handleForgotPassword">
        <div>
          <label for="email" class="block text-sm/6 font-medium text-gray-100">Email address</label>
          <div class="mt-2">
            <Input 
              type="email" 
              name="email" 
              id="email" 
              autocomplete="email" 
              required 
              v-model="email"
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
            :disabled="authStore.isLoading"
            class="flex w-full justify-center rounded-md bg-indigo-500 px-3 py-1.5 text-sm/6 font-semibold text-white hover:bg-indigo-400 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-500"
          >
            {{ authStore.isLoading ? 'Sending...' : 'Send reset link' }}
          </Button>
        </div>
      </form>

      <p class="mt-10 text-center text-sm/6 text-gray-400">
        Remember your password?
        {{ ' ' }}
        <router-link to="/login" class="font-semibold text-indigo-400 hover:text-indigo-300">Sign in here</router-link>
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useAuthStore } from '@/stores/auth';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';

const email = ref('');
const authStore = useAuthStore();

const handleForgotPassword = async () => {
  await authStore.forgotPassword(email.value);
};
</script>

<style scoped>
</style>
