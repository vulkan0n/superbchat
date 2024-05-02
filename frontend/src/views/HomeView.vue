<script>
import { RouterLink } from "vue-router";
export default {
  setup() {
    return { RouterLink };
  },
};
</script>

<template>
  <div class="flex flex-col h-screen justify-center items-center">
    <div class="max-w-xl mx-4 space-y-10 text-gray-600">
      <div class="max-w-xl mx-4 space-y-10 text-gray-600">
        <div class="relative w-full">
          <div
            class="rounded-lg shadow-lg p-8 bg-white m-auto text-center w-full overflow-hidden overflow-ellipsis"
          >
            <b>What is Superbchat?</b><br /><br />
            <p>Use this web app to receive BCH superchats in your stream.</p>
            <p>
              <RouterLink to="/user/signup" class="text-green-500 underline">Sign up</RouterLink> to create your superbchat and start receiving donations,
              or <RouterLink to="/user/login" class="text-green-500 underline">Login</RouterLink> to access your received messages and get your alert URL
              to set up in OBS.
            </p>
            <p>
              You can check the source code
              <a
                href="https://github.com/vulkan0n/superbchat"
                class="text-green-500 underline"
                >Here</a
              >.
            </p>
          </div>
        </div>
      </div>
      <div
        class="flex flex-wrap flex-col sm:flex-row sm:space-x-4 items-center justify-center pb-8 sm:pb-0"
      >
        <div>
          <a
            href="https://discover.cash/"
            target="_blank"
            class="text-green-500 underline"
            >Why Bitcoin Cash?</a
          >
        </div>
        <div>
          <RouterLink to="/user/signup" class="text-green-500 underline">Create a superbchat page</RouterLink>
        </div>
      </div>
    </div>
  </div>
</template>
