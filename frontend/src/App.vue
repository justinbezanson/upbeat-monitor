<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'; // Import the auth store
import { useRouter } from 'vue-router';
import { Button } from '@/components/ui/button'; // Import Button for logout

import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'

import { Settings, LogOut } from 'lucide-vue-next'; 

const authStore = useAuthStore();
const router = useRouter();

const handleLogout = async () => {
  const success = await authStore.logout();
  if (success) {
    router.push('/login'); // Redirect to login page after logout
  }
};
</script>

<template>
  <div v-if="authStore.isAuthenticated" class="grid grid-cols-[250px_1fr_250px] min-h-screen bg-gray-900">

    <div class="flex flex-col p-4 bg-gray-800">
      <h2 class="font-bold text-white mb-4">Left Menu</h2>
      <div class="justify-start w-full text-gray-400">Content at the top</div>
      <div class="mt-auto">
        <DropdownMenu>
          <DropdownMenuTrigger as-child>
            <Button variant="ghost" class="w-full justify-start bg-white">With Icon</Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent side="top" align="start" class="w-[218px]">
            <DropdownMenuItem class="cursor-pointer">
              <Settings class="mr-2 size-4" />
              <span>Settings</span>
            </DropdownMenuItem>
            <DropdownMenuItem @click="handleLogout" class="cursor-pointer">
              <LogOut class="mr-2 size-4" />
              <span>Logout</span>
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </div>

    <div class="p-4">
      <router-view />
    </div>

    <div class="p-4">
      <h3 class="font-bold">Column 3 (1/4)</h3>
      <div class="bg-slate-800 text-gray-400 p-4 mt-4 rounded">This content is in the final quarter-width column.</div>
    </div>

  </div>
  <div v-else class="flex flex-col items-center justify-center min-h-screen">
    <router-view />
  </div>
</template>

<style scoped>
/* Scoped styles if any */
</style>
