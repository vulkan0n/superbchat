<script>
import { RouterLink, RouterView } from "vue-router";
export default {
  setup() {
    return { RouterLink, RouterView };
  },
};
</script>

<template>
  <div
    class="bg-gray-50 flex flex-col min-h-screen"
    style="font-family: Lato, sans-serif"
  >
    <header class="max-w-6xl w-full flex justify-end m-auto mt-4 pr-4">
      <RouterLink to="/">
        <div class="title">
          <img
            alt="BCH logo"
            class="logo"
            src="@/assets/bitcoin-cash-bch-logo.svg"
            width="45"
            height="45"
          />
          <h1>Superbchat!</h1>
        </div>
      </RouterLink>
    </header>
    <div class="mt-2 flex-1 flex items-center justify-center flex-col">
      <RouterView />
    </div>
  </div>
</template>
