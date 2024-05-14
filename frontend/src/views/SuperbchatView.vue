<script>
import { useRouter, useRoute } from "vue-router";
import { onMounted, ref } from "vue";
import { BaseWallet, Wallet } from "mainnet-js";

const fakeUsers = ["vulkan0n", "pepe"];
const walletAddress = ref();
const cashAddress = ref();
const isValidUser = ref(true);

export default {
  setup() {
    const donatorName = ref("Anonymous");
    const donationAmount = ref(0.00000547);
    const donationMessage = ref("");
    const showAmount = ref(true);
    const isCashAddrs = ref(false);

    const donationLowError = ref(false);
    const errorClass = ref("border-red-500");
    const selectedTabClass = ref(
      "inline-block p-4 text-blue-600 bg-gray-100 rounded-t-lg active dark:bg-green-700 dark:text-white"
    );
    const unselectedTabClass = ref(
      "inline-block p-4 rounded-t-lg hover:text-gray-600 hover:bg-gray-50 dark:text-gray-700 dark:hover:bg-gray-800 dark:hover:text-gray-300"
    );

    const router = useRouter();
    const user = useRoute().params.user;

    function verifyAndPay() {
      donationLowError.value = donationAmount.value < 0.00000547;

      if (!donationLowError.value) {
        console.log(
          `Donator: ${donatorName.value} - Message: ${donationMessage.value} 
          Amount: ${donationAmount.value}`
        );
      }
    }

    function onQRChange() {
      const currentVal = isCashAddrs.value;
      isCashAddrs.value = !currentVal;
    }

    onMounted(async () => {
      if (!fakeUsers.includes(user)) {
        router.push({ name: "404" });
        isValidUser.value = false;
      } else {
        BaseWallet.StorageProvider = IndexedDBProvider;
        const wallet = await Wallet.named(`user:${user}`);
        console.log(wallet);
        walletAddress.value = wallet.address;
        cashAddress.value = wallet.tokenaddr;
      }
    });

    return {
      donatorName,
      donationAmount,
      donationMessage,
      errorClass,
      selectedTabClass,
      unselectedTabClass,
      showAmount,
      donationLowError,
      user,
      walletAddress,
      isValidUser,
      verifyAndPay,
      isCashAddrs,
      cashAddress,
      onQRChange,
    };
  },
};
</script>

<template>
  <div class="flex flex-col h-screen items-center">
    <div class="max-w-lg mx-4 space-y-10 text-gray-600">
      <div class="max-w-lg mx-4 space-y-10 text-gray-600">
        <div class="relative w-full">
          <div
            class="rounded-lg shadow-lg p-8 bg-white m-auto w-full overflow-hidden overflow-ellipsis"
          >
            <h1 class="block text-gray-700 text-lg font-bold mb-3">
              {{ user.substring(0, 1).toUpperCase() + user.substring(1) }}'s
              Superbchat
            </h1>
            <hr class="mb-4" />
            <div class="mb-4">
              <label
                class="block text-gray-700 text-sm font-bold mb-2"
                for="donatorName"
              >
                Name:
              </label>
              <input
                class="shadow appearance-none border rounded w-80 py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
                id="donatorName"
                type="text"
                placeholder="Anonymous"
                v-model="donatorName"
              />
            </div>
            <div class="mb-4">
              <label
                class="block text-gray-700 text-sm font-bold mb-2"
                for="donationAmount"
              >
                Donation Amount:
              </label>
              <input
                class="shadow appearance-none border rounded w-80 py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
                :class="donationLowError ? errorClass : ''"
                id="donationAmount"
                type="number"
                v-model="donationAmount"
              />
              <p v-if="donationLowError" class="text-red-500 text-xs italic">
                Minimun donation amount must be higher than 0.00000546.
              </p>
            </div>
            <div class="mb-4">
              <label
                class="block text-gray-700 text-sm font-bold mb-2"
                for="donationMessage"
              >
                Message:
              </label>
              <input
                class="shadow appearance-none border rounded w-80 py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
                id="donationMessage"
                type="text"
                placeholder=""
                v-model="donationMessage"
              />
            </div>
            <div class="mb-6">
              <label class="block text-gray-500 font-bold">
                <input
                  class="mr-2 leading-tight"
                  type="checkbox"
                  v-model="showAmount"
                />
                <span class="text-sm"> Show amount on stream? </span>
              </label>
            </div>
            <ul
              class="flex flex-wrap text-sm font-medium text-center text-gray-500 border-b border-gray-200 dark:border-gray-700 dark:text-gray-400"
            >
              <li class="me-2">
                <a
                  href="#"
                  aria-current="page"
                  :class="isCashAddrs ? unselectedTabClass : selectedTabClass"
                  @click="onQRChange"
                  >Address</a
                >
              </li>
              <li class="me-2">
                <a
                  href="#"
                  :class="isCashAddrs ? selectedTabClass : unselectedTabClass"
                  @click="onQRChange"
                  >Cash Address</a
                >
              </li>
            </ul>

            <div class="w-80" v-show="isValidUser">
              <span class="text-xs">{{!isCashAddrs ? walletAddress : cashAddress}}</span>
              <qr-code
                id="qr-bch"
                :contents="!isCashAddrs ? walletAddress : cashAddress"
                module-color="#1c7d43"
                position-ring-color="#13532d"
                position-center-color="#70c559"
                style="background-color: #fff"
                class="qr"
              />
            </div>
            <div class="flex items-center justify-between">
              <button
                class="bg-green-600 hover:bg-green-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
                type="button"
                @click="verifyAndPay"
              >
                Send Donation!
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
