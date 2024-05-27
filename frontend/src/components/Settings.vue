<script>
import { ref } from "vue";
import { RouterLink, useRouter } from "vue-router";

export default {
  setup() {
    const address = ref("bitcoincash:addressTest");
    const minDonation = ref(0.00000547);
    const defaultShowAmount = ref(true);

    const minDonationLowError = ref(false);
    const errorClass = ref("border-red-500");

    const router = useRouter();

    function verifyAndUpdateSettings() {
      minDonationLowError.value = minDonation.value < 0.00000547;

      if (!minDonationLowError.value) {
        router.push("/dashboard");
      }
    }

    return {
      RouterLink,
      address,
      minDonation,
      defaultShowAmount,
      minDonationLowError,
      errorClass,
      verifyAndUpdateSettings,
    };
  },
};
</script>

<template>
  <div class="flex flex-col  justify-center items-center">
    <div class="max-w-lg mx-4 space-y-10 text-gray-600">
      <div class="max-w-lg mx-4 space-y-10 text-gray-600">
        <div class="relative w-full">
          <div
            class="rounded-lg shadow-lg p-8 bg-white m-auto w-full overflow-hidden overflow-ellipsis"
          >
            <div class="mb-4">
              <label
                class="block text-gray-700 text-sm font-bold mb-2"
                for="address"
              >
                Address
              </label>
              <input
                class="shadow appearance-none border rounded w-80 py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
                id="address"
                type="text"
                placeholder="bitcoincash:"
                v-model="address"
              />
            </div>
            <div class="mb-4">
              <label
                class="block text-gray-700 text-sm font-bold mb-2"
                for="minDonation"
              >
                Minimun Donation Amount
              </label>
              <input
                class="shadow appearance-none border rounded w-80 py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
                :class="minDonationLowError ? errorClass : ''"
                id="minDonation"
                type="number"
                v-model="minDonation"
              />
              <p v-if="minDonationLowError" class="text-red-500 text-xs italic">
                Minimun donation amount must be higher than 0.00000546.
              </p>
            </div>
            <div class="mb-6">
              <label class="block text-gray-500 font-bold">
                <input class="mr-2 leading-tight accent-emerald-700" type="checkbox" v-model="defaultShowAmount" />
                <span class="text-sm"> Default "Show Amount" checked </span>
              </label>
            </div>
            <div class="flex items-center justify-between">
              <button
                class="bg-emerald-600 hover:bg-emerald-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
                type="button"
                @click="verifyAndUpdateSettings"
              >
                Update
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
