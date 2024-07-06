<script>
import { ref } from "vue";
import Dashboard from "../components/Dashboard.vue"
import Settings from "@/components/Settings.vue";
import Wallet from "@/components/Wallet.vue";

const fakeUser = "vulkan0n";

export default {
  components: { Dashboard, Settings, Wallet },
  setup() {
    let showNavMenuMobile = ref(false);
    let selectedTab = ref("Dashboard");
    return { Dashboard, fakeUser, showNavMenuMobile, selectedTab };
  },
};
</script>

<template>
  <nav
    class="bg-white dark:bg-gray-700 fixed w-full z-20 top-0 start-0 border-b border-gray-200 dark:border-gray-600"
  >
    <div
      class="max-w-screen-xl flex flex-wrap items-center justify-between mx-auto p-4"
    >
      <div class="flex items-center space-x-3 rtl:space-x-reverse">
        <img
          src="../../src/assets/bitcoin-cash-bch-logo.svg"
          class="h-8"
          alt="BCH Logo"
        />
        <span
          class="self-center text-2xl font-semibold whitespace-nowrap dark:text-white"
          >Superbchat</span
        >
      </div>
      <div class="flex md:order-2 space-x-3 md:space-x-0 rtl:space-x-reverse">
        <RouterLink
          :to="'/' + fakeUser"
          target="_blank"
          class="bg-emerald-500 hover:bg-emerald-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
          >Your page</RouterLink
        >
        <button
          data-collapse-toggle="navbar-sticky"
          type="button"
          class="inline-flex items-center p-2 w-10 h-10 justify-center text-sm text-gray-500 rounded-lg md:hidden hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-gray-200 dark:text-gray-400 dark:hover:bg-gray-700 dark:focus:ring-gray-600"
          aria-controls="navbar-sticky"
          aria-expanded="false"
          @click="showNavMenuMobile = !showNavMenuMobile"
        >
          <span class="sr-only">Open main menu</span>
          <svg
            class="w-5 h-5"
            aria-hidden="true"
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 17 14"
          >
            <path
              stroke="currentColor"
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M1 1h15M1 7h15M1 13h15"
            />
          </svg>
        </button>
      </div>
      <div
        class="items-center justify-between w-full md:flex md:w-auto md:order-1"
        :class="showNavMenuMobile ? '' : 'hidden'"
        id="navbar-sticky"
      >
        <ul
          class="flex flex-col p-4 md:p-0 mt-4 font-medium border border-gray-100 rounded-lg bg-gray-50 md:space-x-8 rtl:space-x-reverse md:flex-row md:mt-0 md:border-0 md:bg-white dark:bg-gray-800 md:dark:bg-gray-700 dark:border-gray-700"
        >
          <li>
            <button
              class="block py-2 px-3 text-white rounded md:p-0 md:text-blue-700 md:dark:text-white"
              :class="selectedTab == 'Dashboard' ? 'font-bold bg-emerald-700 md:bg-transparent' : 'md:hover:underline md:dark:hover:underline dark:hover:bg-gray-700 '"
              @click="selectedTab = 'Dashboard'"
              >Dashboard</button
            >
          </li>
          <li>
            <button
              class="block py-2 px-3 text-gray-900 rounded md:p-0 md:hover:bg-transparent dark:text-white dark:hover:text-white md:dark:hover:bg-transparent dark:border-gray-700"
              :class="selectedTab == 'Settings' ? 'font-bold bg-emerald-700 md:bg-transparent' : 'md:hover:underline md:dark:hover:underline dark:hover:bg-gray-700 '"
              @click="selectedTab = 'Settings'"
              >Settings</button
            >
          </li>
          <li>
            <button
              class="block py-2 px-3 text-gray-900 rounded md:p-0 md:hover:bg-transparent  dark:text-white dark:hover:text-white  dark:border-gray-700"
              :class="selectedTab == 'Wallet' ? 'font-bold bg-emerald-700 md:bg-transparent' : 'md:hover:underline md:dark:hover:underline dark:hover:bg-gray-700 '"
              @click="selectedTab = 'Wallet'"
              >Wallet</button
            >
          </li>
        </ul>
      </div>
    </div>
  </nav>
  <div class="mt-24 mb-14 mx-3 md:mx-20">
    <dashboard v-show="selectedTab == 'Dashboard'" />
    <settings v-show="selectedTab == 'Settings'" />
    <wallet v-show="selectedTab == 'Wallet'" />
  </div>
</template>
