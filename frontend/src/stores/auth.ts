import { defineStore } from 'pinia';
import { ref, computed } from 'vue'; // Import computed

export const useAuthStore = defineStore('auth', () => {
  const isLoading = ref(false);
  const error = ref<string | null>(null);
  const successMessage = ref<string | null>(null);
  const userToken = ref<string | null>(localStorage.getItem('userToken')); // Initialize from localStorage

  const isAuthenticated = computed(() => !!userToken.value); // New getter

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

  const login = async (email: string, password: string) => {
    isLoading.value = true;
    error.value = null;
    successMessage.value = null; // Clear success message from registration

    try {
      const response = await fetch('/api/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email, password }),
      });

      const data = await response.json();

      if (response.ok && data.token) {
        userToken.value = data.token;
        localStorage.setItem('userToken', data.token); // Persist token
        successMessage.value = data.message || 'Login successful!';
        return true;
      } else {
        error.value = data.error || 'Login failed.';
        return false;
      }
    } catch (e: any) {
      error.value = 'An unexpected error occurred: ' + e.message;
      return false;
    } finally {
      isLoading.value = false;
    }
  };

  const logout = async () => {
    isLoading.value = true;
    error.value = null;
    successMessage.value = null; // Clear messages on logout

    try {
      // Optionally call backend logout endpoint (e.g., for token invalidation on server)
      // For simple JWTs, client-side token removal is often sufficient.
      await fetch('/api/logout', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${userToken.value}`,
        },
      });

      userToken.value = null;
      localStorage.removeItem('userToken');
      successMessage.value = 'Logged out successfully.';
      return true;
    } catch (e: any) {
      error.value = 'An unexpected error occurred during logout: ' + e.message;
      return false;
    } finally {
      isLoading.value = false;
    }
  };

  return {
    isLoading,
    error,
    successMessage,
    userToken,
    isAuthenticated, // Expose isAuthenticated
    register,
    login,
    logout, // Expose logout action
  };
});