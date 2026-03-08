import { defineStore } from 'pinia';
import { ref, computed } from 'vue'; // Import computed

export const useAuthStore = defineStore('auth', () => {
  const isLoading = ref(false);
  const error = ref<string | null>(null);
  const successMessage = ref<string | null>(null);
  const userToken = ref<string | null>(null); 

  const _isAuthenticated = ref(false); 
  const isAuthenticated = computed(() => _isAuthenticated.value);
  const isInitialized = ref(false); // New flag

  const initializeAuth = async () => {
    isLoading.value = true;
    try {
      const response = await fetch('/api/ping', {
        method: 'GET',
        credentials: 'include',
      });
      _isAuthenticated.value = response.ok;
    } catch (e) {
      _isAuthenticated.value = false;
    } finally {
      isLoading.value = false;
      isInitialized.value = true; // Mark as done
    }
  };

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
        credentials: 'include',
        body: JSON.stringify({ email, password }),
      });

      const data = await response.json();

      if (response.ok) {
        successMessage.value = data.message || 'Registration successful!';
        _isAuthenticated.value = true;
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
    successMessage.value = null;

    try {
      const response = await fetch('/api/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({ email, password }),
      });

      const data = await response.json();

      if (response.ok) {
        _isAuthenticated.value = true;
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
    successMessage.value = null;

    try {
      await fetch('/api/logout', {
        method: 'POST',
        credentials: 'include',
      });

      successMessage.value = 'Logged out successfully.';
      return true;
    } catch (e: any) {
      error.value = 'An unexpected error occurred during logout: ' + e.message;
      return false;
    } finally {
      _isAuthenticated.value = false;
      isLoading.value = false;
    }
  };

  const forgotPassword = async (email: string) => {
    isLoading.value = true;
    error.value = null;
    successMessage.value = null;

    try {
      const response = await fetch('/api/forgot-password', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({ email }),
      });

      const data = await response.json();

      if (response.ok) {
        successMessage.value = data.message;
        return true;
      } else {
        error.value = data.error || 'Failed to send reset link.';
        return false;
      }
    } catch (e: any) {
      error.value = 'An unexpected error occurred: ' + e.message;
      return false;
    } finally {
      isLoading.value = false;
    }
  };

  const resetPassword = async (token: string, password: string) => {
    isLoading.value = true;
    error.value = null;
    successMessage.value = null;

    try {
      const response = await fetch('/api/reset-password', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({ token, password }),
      });

      const data = await response.json();

      if (response.ok) {
        successMessage.value = data.message;
        return true;
      } else {
        error.value = data.error || 'Failed to reset password.';
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
          userToken,
          isAuthenticated, 
          isInitialized, // Expose isInitialized
          initializeAuth, 
          register,
          login,
          logout, 
          forgotPassword,
          resetPassword,
        };
      });
      