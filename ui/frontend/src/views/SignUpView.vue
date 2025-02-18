<script>
import { ref } from "vue";
import { RouterLink, useRouter } from "vue-router";
import axios from "axios";

const invalidUsernames = [
  "user-signup",
  "user-login",
  "verify-tkn",
  "user-id",
  "user",
  "settings",
  "superbchat",
  "superbchat-get",
  "login",
  "signup",
  "dashboard",
  "privacy-policy",
  "404",
];

// Must start with a letter.
// Can contain letters, numbers, and underscores (_).
// Length: 3 to 15 characters.
const usernameRegex = /^[a-zA-Z][a-zA-Z0-9_]{2,14}$/;

export default {
  setup() {
    const username = ref("");
    const password = ref("");
    const repeatedPassword = ref("");

    const emptyPasswordError = ref(false);
    const emptyUserError = ref(false);
    const shortPasswordError = ref(false);
    const differentPasswordError = ref(false);
    const usernameTakenError = ref(false);
    const usernameInvalidError = ref(false);
    const isLoading = ref(false);

    const errorClass = ref("border-red-500");

    const router = useRouter();

    function verifyAndSignUp() {
      emptyUserError.value = username.value == "";
      emptyPasswordError.value = password.value == "";
      shortPasswordError.value = password.value.length < 8;
      differentPasswordError.value = password.value != repeatedPassword.value;

      if (invalidUsernames.includes(username.value)) {
        usernameTakenError.value = true;
        setTimeout(() => {
          usernameTakenError.value = false;
        }, 3000);
      }

      if (!usernameRegex.test(username.value)) {
        usernameInvalidError.value = true;
        setTimeout(() => {
          usernameInvalidError.value = false;
        }, 3000);
      }

      if (
        !emptyUserError.value &&
        !emptyPasswordError.value &&
        !shortPasswordError.value &&
        !differentPasswordError.value &&
        !usernameTakenError.value &&
        !usernameInvalidError.value
      ) {
        postUserSignUp();
      }
    }

    async function postUserSignUp() {
      isLoading.value = true;

      try {
        const response = await axios.post("/user-signup", {
          user: username.value,
          pass: password.value,
        });

        if (response.statusText == "OK") {
          isLoading.value = true;
          try {
            const loginResponse = await axios.post("/user-login", {
              user: username.value,
              pass: password.value,
            });
            if (loginResponse.statusText == "OK") {
              const token = loginResponse.data.token;
              const userId = loginResponse.data.userId;
              localStorage.setItem("token", token);
              localStorage.setItem("userId", userId);
            }
          } catch (err) {
            console.log(err);
          } finally {
            router.push("/dashboard");
          }
        } else {
          console.log(response);
          usernameTakenError.value = true;
          setTimeout(() => {
            usernameTakenError.value = false;
          }, 3000);
        }
      } catch (err) {
        console.log(err);
        usernameTakenError.value = true;
        setTimeout(() => {
          usernameTakenError.value = false;
        }, 3000);
      } finally {
        isLoading.value = false;
      }
    }

    return {
      RouterLink,
      username,
      password,
      repeatedPassword,
      emptyUserError,
      emptyPasswordError,
      shortPasswordError,
      differentPasswordError,
      usernameTakenError,
      usernameInvalidError,
      errorClass,
      verifyAndSignUp,
      isLoading,
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
                :class="
                  emptyUserError || usernameTakenError || usernameInvalidError
                    ? errorClass
                    : ''
                "
                id="username"
                type="text"
                placeholder="Username"
                v-model="username"
              />
              <p v-if="emptyUserError" class="text-red-500 text-xs italic">
                Please choose a username.
              </p>
              <p v-if="usernameTakenError" class="text-red-500 text-xs italic">
                This username is already in use.
              </p>
              <p v-if="usernameInvalidError" class="text-red-500 text-xs italic">
                This username is invalid.
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
                class="shadow appearance-none border rounded w-80 py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
                :class="
                  emptyPasswordError || shortPasswordError ? errorClass : ''
                "
                id="password"
                type="password"
                placeholder="•••••••••"
                v-model="password"
              />
              <p v-if="emptyPasswordError" class="text-red-500 text-xs italic">
                Please choose a password.
              </p>
              <p v-if="shortPasswordError" class="text-red-500 text-xs italic">
                Password must be at least 8 characters long.
              </p>
            </div>
            <div class="mb-4">
              <label
                class="block text-gray-700 text-sm font-bold mb-2"
                for="repeatedPassword"
              >
                Repeat Password
              </label>
              <input
                class="shadow appearance-none border rounded w-80 py-2 px-3 text-gray-700 mb-3 leading-tight focus:outline-none focus:shadow-outline"
                :class="differentPasswordError ? errorClass : ''"
                id="repeatedPassword"
                type="password"
                placeholder="•••••••••"
                v-model="repeatedPassword"
              />
              <p
                v-if="differentPasswordError"
                class="text-red-500 text-xs italic"
              >
                Password does not match.
              </p>
            </div>

            <div class="flex items-center justify-between">
              <button
                class="bg-green-600 hover:bg-green-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
                type="button"
                @click="verifyAndSignUp"
              >
                {{ isLoading ? "Loading..." : "Create" }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
