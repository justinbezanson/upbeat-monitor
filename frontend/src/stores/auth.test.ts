import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { setActivePinia, createPinia } from 'pinia';
import { useAuthStore } from './auth';

describe('Auth Store', () => {
  let originalFetch: typeof fetch; // Change global.fetch to fetch

  beforeEach(() => {
    setActivePinia(createPinia());
    vi.clearAllMocks();
    localStorage.clear(); // Clear localStorage before each test
    originalFetch = fetch; // Store original fetch
  });

  afterEach(() => {
    fetch = originalFetch; // Restore original fetch
  });

  // Mock fetch API helper
  const mockFetch = (response: any, ok: boolean = true, status: number = 200) => {
    return vi.fn(() =>
      Promise.resolve({
        ok,
        status,
        json: () => Promise.resolve(response),
      }),
    ) as any;
  };

  describe('register', () => {
    it('should register a user successfully', async () => {
      const authStore = useAuthStore();
      const successResponse = { message: 'Registration successful!', user_id: '123' };
      fetch = mockFetch(successResponse); // Change global.fetch to fetch

      const success = await authStore.register('test@example.com', 'password123');

      expect(success).toBe(true);
      expect(authStore.isLoading).toBe(false);
      expect(authStore.error).toBeNull();
      expect(authStore.successMessage).toBe('Registration successful!');
      expect(fetch).toHaveBeenCalledWith('/api/register', expect.any(Object)); // Change global.fetch to fetch
    });

    it('should handle registration failure', async () => {
      const authStore = useAuthStore();
      const errorResponse = { error: 'User already exists' };
      fetch = mockFetch(errorResponse, false, 409); // Change global.fetch to fetch

      const success = await authStore.register('test@example.com', 'password123');

      expect(success).toBe(false);
      expect(authStore.isLoading).toBe(false);
      expect(authStore.error).toBe('User already exists');
      expect(authStore.successMessage).toBeNull();
    });

    it('should handle network errors during registration', async () => {
      const authStore = useAuthStore();
      fetch = vi.fn(() => Promise.reject(new Error('Network Error'))) as any; // Change global.fetch to fetch

      const success = await authStore.register('test@example.com', 'password123');

      expect(success).toBe(false);
      expect(authStore.isLoading).toBe(false);
      expect(authStore.error).toContain('Network Error');
      expect(authStore.successMessage).toBeNull();
    });
  });

  describe('login', () => {
    it('should log in a user successfully and store token', async () => {
      const authStore = useAuthStore();
      const successResponse = { message: 'Logged in successfully', token: 'fake-jwt-token', user_id: '123' };
      fetch = mockFetch(successResponse); // Change global.fetch to fetch

      const success = await authStore.login('test@example.com', 'password123');

      expect(success).toBe(true);
      expect(authStore.isLoading).toBe(false);
      expect(authStore.error).toBeNull();
      expect(authStore.successMessage).toBe('Logged in successfully');
      expect(authStore.userToken).toBe('fake-jwt-token');
      expect(localStorage.getItem('userToken')).toBe('fake-jwt-token');
      expect(authStore.isAuthenticated).toBe(true);
      expect(fetch).toHaveBeenCalledWith('/api/login', expect.any(Object)); // Change global.fetch to fetch
    });

    it('should handle login failure', async () => {
      const authStore = useAuthStore();
      const errorResponse = { error: 'Invalid credentials' };
      fetch = mockFetch(errorResponse, false, 401); // Change global.fetch to fetch

      const success = await authStore.login('test@example.com', 'password123');

      expect(success).toBe(false);
      expect(authStore.isLoading).toBe(false);
      expect(authStore.error).toBe('Invalid credentials');
      expect(authStore.successMessage).toBeNull();
      expect(authStore.userToken).toBeNull();
      expect(localStorage.getItem('userToken')).toBeNull();
      expect(authStore.isAuthenticated).toBe(false);
    });

    it('should handle network errors during login', async () => {
      const authStore = useAuthStore();
      fetch = vi.fn(() => Promise.reject(new Error('Network Error'))) as any; // Change global.fetch to fetch

      const success = await authStore.login('test@example.com', 'password123');

      expect(success).toBe(false);
      expect(authStore.isLoading).toBe(false);
      expect(authStore.error).toContain('Network Error');
      expect(authStore.successMessage).toBeNull();
      expect(authStore.userToken).toBeNull();
      expect(localStorage.getItem('userToken')).toBeNull();
      expect(authStore.isAuthenticated).toBe(false);
    });
  });

  describe('logout', () => {
    it('should log out a user successfully and clear token', async () => {
      const authStore = useAuthStore();
      authStore.userToken = 'existing-jwt-token'; // Simulate logged in
      localStorage.setItem('userToken', 'existing-jwt-token');
      fetch = mockFetch({ message: 'Logged out successfully' }); // Change global.fetch to fetch

      const success = await authStore.logout();

      expect(success).toBe(true);
      expect(authStore.isLoading).toBe(false);
      expect(authStore.error).toBeNull();
      expect(authStore.successMessage).toBe('Logged out successfully.');
      expect(authStore.userToken).toBeNull();
      expect(localStorage.getItem('userToken')).toBeNull();
      expect(authStore.isAuthenticated).toBe(false);
      expect(fetch).toHaveBeenCalledWith('/api/logout', expect.any(Object)); // Change global.fetch to fetch
    });

    it('should handle logout when not authenticated (no token)', async () => {
        const authStore = useAuthStore();
        // userToken is already null from beforeEach
        fetch = mockFetch({ message: 'Logged out successfully' }); // Mock response for the fetch call inside logout
        const success = await authStore.logout();
        expect(success).toBe(true); // Should still return true, as there's nothing to clear
        expect(authStore.userToken).toBeNull();
        expect(localStorage.getItem('userToken')).toBeNull();
    });

    it('should handle network errors during logout', async () => {
      const authStore = useAuthStore();
      authStore.userToken = 'existing-jwt-token';
      localStorage.setItem('userToken', 'existing-jwt-token'); // Simulate logged in
      fetch = vi.fn(() => Promise.reject(new Error('Network Error'))) as any; // Change global.fetch to fetch

      const success = await authStore.logout();

      expect(success).toBe(false);
      expect(authStore.isLoading).toBe(false);
      expect(authStore.error).toContain('Network Error');
      expect(authStore.successMessage).toBeNull();
      // Even if network error, token should be cleared client-side
      expect(authStore.userToken).toBeNull();
      expect(localStorage.getItem('userToken')).toBeNull();
      expect(authStore.isAuthenticated).toBe(false);
    });
  });

  describe('isAuthenticated getter', () => {
    it('should return true if userToken is present', () => {
      const authStore = useAuthStore();
      authStore.userToken = 'some-token';
      expect(authStore.isAuthenticated).toBe(true);
    });

    it('should return false if userToken is null', () => {
      const authStore = useAuthStore();
      authStore.userToken = null;
      expect(authStore.isAuthenticated).toBe(false);
    });
  });
});
