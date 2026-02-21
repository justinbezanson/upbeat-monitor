import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useAuthStore = defineStore('auth', () => {
  const isLoading = ref(false);
  const error = ref<string | null>(null);
  const successMessage = ref<string | null>(null);

  const register = async (email: string, password: string) => {
    isLoading.value = true;
    error.value = null;
    successMessage.value = null;

    try {
      const response = await fetch('/api/register', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email, password }),
      });

      const data = await response.json();

      if (response.ok) {
        successMessage.value = data.message || 'Registration successful!';
        return true;
      } else {
        error.value = data.error || 'Registration failed.';
        return false;
      }
    } catch (e: any) {
      error.value = 'An unexpected error occurred: ' + e.message;
      return false;
    } finally {
      isLoading.value = false;
    }
  };

  return {
    isLoading,
    error,
    successMessage,
    register,
  };
});