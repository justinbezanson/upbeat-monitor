<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'; // Import the auth store
import { useRouter } from 'vue-router';
import { Button } from '@/components/ui/button'; // Import Button for logout

import {
  NavigationMenu,
  NavigationMenuContent,
  NavigationMenuItem,
  NavigationMenuLink,
  NavigationMenuList,
  NavigationMenuTrigger,
} from '@/components/ui/navigation-menu'

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

    <div class="p-4 bg-gray-800">
      <h2 class="font-bold text-white mb-4">Left Menu</h2>
      <div class="flex flex-col gap-2">
        <Button @click="handleLogout" variant="ghost" class="bg-sky-600 text-sky-400 hover:bg-sky-700 justify-start w-full">Logout</Button>
        <NavigationMenu :disableHoverTrigger="true">
          <NavigationMenuList>
            <NavigationMenuItem>
              <NavigationMenuTrigger>With Icon</NavigationMenuTrigger>
              <NavigationMenuContent>
                <ul class="grid w-[200px] gap-1 p-2 bg-popover rounded-md">
                  <li>
                    <NavigationMenuLink as-child>
                      <a href="#" class="flex items-center gap-2 p-2 hover:bg-accent hover:text-accent-foreground rounded-sm transition-colors">
                        <Settings class="size-4" />
                        <span>Settings</span>
                      </a>
                    </NavigationMenuLink>
                  </li>
                  <li>
                    <NavigationMenuLink as-child>
                      <a href="#" class="flex items-center gap-2 p-2 hover:bg-accent hover:text-accent-foreground rounded-sm transition-colors">
                        <LogOut class="size-4" />
                        <span>Logout</span>
                      </a>
                    </NavigationMenuLink>
                  </li>
                </ul>
              </NavigationMenuContent>
            </NavigationMenuItem>
          </NavigationMenuList>
        </NavigationMenu>
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

