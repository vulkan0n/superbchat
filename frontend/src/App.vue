<script>
import { RouterLink, RouterView, useRoute } from "vue-router";
import { ref } from "vue";

export default {
  setup() {
    const route = useRoute();
    const routesWithHeader = ["home", "superbchat", "404"];
    let show = ref(false);
    return {
      RouterLink,
      RouterView,
      route,
      routesWithHeader,
      show,
    };
  },
};
</script>

<template>
  <div class="bg-gray-50 flex flex-col min-h-screen">
    <header
      v-show="routesWithHeader.includes(route.name)"
      class="max-w-6xl w-full flex justify-end m-auto mt-4 pr-4"
    >
      <div>
        <div class="relative">
          <!-- Dropdown toggle button -->
          <button
            @click="show = !show"
            class="flex items-center p-2 text-white bg-emerald-600 rounded-md font-bold"
          >
            <span class="mr-2">Get Started</span>
            <svg
              class="w-5 h-5 text-indigo-100 dark:text-white"
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 20 20"
              fill="currentColor"
            >
              <path
                fill-rule="evenodd"
                d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z"
                clip-rule="evenodd"
              />
            </svg>
          </button>

          <!-- Dropdown menu -->
          <div
            v-show="show"
            class="absolute right-0 py-2 mt-2 bg-emerald-500 rounded-md shadow-xl w-44"
          >
            <router-link
              to="/"
              class="block px-4 py-2 text-sm text-white hover:bg-emerald-400"
              v-show="route.name != 'home'"
            >
              What is Superbchat?
            </router-link>
            <router-link
              to="/user/signup"
              class="block px-4 py-2 text-sm text-white hover:bg-emerald-400"
            >
              Create a page
            </router-link>
            <router-link
              to="/user/login"
              class="block px-4 py-2 text-sm text-white hover:bg-emerald-400"
            >
              Login
            </router-link>
          </div>
        </div>
      </div>
    </header>
    <div class="mt-2 flex-1 flex items-center justify-center flex-col">
      <RouterView />
    </div>
  </div>
</template>
