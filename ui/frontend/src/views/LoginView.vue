<script>
import { ref } from "vue";
import { RouterLink, useRouter } from "vue-router";

const fakeCredentials = { user: "vulkan0n", password: "test" };

export default {
  setup() {
    const username = ref("");
    const password = ref("");

    const emptyPasswordError = ref(false);
    const emptyUserError = ref(false);
    const credentialsError = ref(false);
    const errorClass = ref("border-red-500");

    const router = useRouter();

    function verifyAndLogin() {
      emptyUserError.value = username.value == "";
      emptyPasswordError.value = password.value == "";

      if (!emptyUserError.value && !emptyPasswordError.value) {
        emptyUserError.value = false;
        emptyPasswordError.value = false;

        if (
          username.value != fakeCredentials.user ||
          password.value != fakeCredentials.password
        ) {
          credentialsError.value = true;
        } else {
          router.push('/dashboard')
        }
      }
    }

    return {
      RouterLink,
      username,
      password,
      emptyUserError,
      emptyPasswordError,
      credentialsError,
      errorClass,
      verifyAndLogin,
    };
  },
};
</script>

<template>
  <div class="flex flex-col h-screen justify-center items-center">
    <div class="max-w-lg mx-4 space-y-10 text-gray-600">
      <div class="max-w-lg mx-4 space-y-10 text-gray-600">
        <div class="relative w-full">
          <div
            class="rounded-lg shadow-lg p-8 bg-white m-auto w-full overflow-hidden overflow-ellipsis"
          >
            <div class="mb-4">
              <label
                class="block text-gray-700 text-sm font-bold mb-2"
                for="username"
              >
                Username
              </label>
              <input
                class="shadow appearance-none border rounded w-80 py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
                :class="emptyUserError || credentialsError ? errorClass : ''"
                id="username"
                type="text"
                placeholder="Username"
                v-model="username"
              />
              <p v-if="emptyUserError" class="text-red-500 text-xs italic">
                Please choose a username.
              </p>
            </div>
            <div class="mb-4">
              <label
                class="block text-gray-700 text-sm font-bold mb-2"
                for="password"
              >
                Password
              </label>
              <input
                class="shadow appearance-none border rounded w-80 py-2 px-3 text-gray-700 mb-3 leading-tight focus:outline-none focus:shadow-outline"
                :class="
                  emptyPasswordError || credentialsError ? errorClass : ''
                "
                id="password"
                type="password"
                placeholder="•••••••••"
                v-model="password"
              />
              <p v-if="emptyPasswordError" class="text-red-500 text-xs italic">
                Please choose a password.
              </p>
              <p v-if="credentialsError" class="text-red-500 text-xs italic">
                Username or password incorrect.
              </p>
            </div>
            <div class="flex items-center justify-between">
              <button
                class="bg-green-600 hover:bg-green-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
                type="button"
                @click="verifyAndLogin"
              >
                Login
              </button>
              <RouterLink
                to="/user/signup"
                class="inline-block align-baseline font-bold text-sm text-green-600 hover:text-green-800"
                >Create a page</RouterLink
              >
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
