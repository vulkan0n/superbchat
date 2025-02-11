<script>
import { useRouter, useRoute } from "vue-router";
import { onMounted, ref } from "vue";
import axios from "axios";

const walletAddress = ref();
const tknEnabled = ref(true);
const cashAddress = ref();
const isValidUser = ref(true);
const minDonation = ref(0.00000547);
const donationAmount = ref(0.00000547);
const msgMaxChar = ref(300);
const showAmount = ref(true);

export default {
  setup() {
    const donatorName = ref("Anonymous");
    const donationMessage = ref("");
    const isCashAddrs = ref(false);
    var userId = 0;
    const isValidSend = ref(false);

    const donationLowError = ref(false);
    const errorClass = ref("border-red-500");
    const selectedTabClass = ref(
      "inline-block p-4 text-blue-600 bg-gray-100 rounded-t-lg active dark:bg-green-700 dark:text-white"
    );
    const unselectedTabClass = ref(
      "inline-block p-4 rounded-t-lg hover:text-gray-600 hover:bg-gray-50 dark:text-gray-700 dark:hover:bg-gray-800 dark:hover:text-gray-300"
    );

    const router = useRouter();
    const route = useRoute();
    const user = useRoute().params.user;

    async function verifyAndPay() {
      if (donationAmount.value < minDonation.value) {
        donationLowError.value = true;
        setTimeout(() => {
          donationLowError.value = false;
        }, 3000);
      }
      if (!donationLowError.value) {
        try {
          const superbchatResponse = await axios.post("/superbchat", {
            name: donatorName.value,
            message: donationMessage.value,
            amount: donationAmount.value,
            isHidden: !showAmount.value,
            recipient: userId,
            isTkn: isCashAddrs.value,
          });
          if (superbchatResponse.statusText == "OK") {
            console.log("Superchat Sent");
            isValidSend.value = true;
            setTimeout(() => {
              isValidSend.value = false;
            }, 3000);
          }
        } catch (err) {
          console.log(err);
        } finally {
          setTimeout(() => {
            router.push(route.fullPath);
          }, 3000);
        }
      }
    }

    function onQRChange() {
      const currentVal = isCashAddrs.value;
      isCashAddrs.value = !currentVal;
    }

    onMounted(async () => {
      try {
        const userInfoResponse = await axios.get("/user/" + user);
        if (userInfoResponse.statusText == "OK") {
          walletAddress.value = userInfoResponse.data.address;
          tknEnabled.value = userInfoResponse.data.tknsEnabled;
          cashAddress.value = userInfoResponse.data.tknAddress;
          showAmount.value = userInfoResponse.data.showAmount;
          minDonation.value = userInfoResponse.data.minDonation;
          donationAmount.value = userInfoResponse.data.minDonation;
          msgMaxChar.value = userInfoResponse.data.messageMaxChars;
          userId = userInfoResponse.data.userId;
        } else {
          router.push({ name: "404" });
          isValidUser.value = false;
        }
      } catch (err) {
        router.push({ name: "404" });
        isValidUser.value = false;
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
      tknEnabled,
      minDonation,
      msgMaxChar,
      isValidSend,
    };
  },
};
</script>

<template>
  <div class="flex flex-col items-center">
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
                maxlength="15"
                v-model="donatorName"
              />
            </div>
            <div v-if="!isCashAddrs" class="mb-4">
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
                step="0.000001"
                v-model="donationAmount"
              />
              <p v-if="donationLowError" class="text-red-500 text-xs italic">
                Minimun donation amount must be higher than {{ minDonation }}
              </p>
            </div>
            <div class="mb-4">
              <label
                class="block text-gray-700 text-sm font-bold mb-2"
                for="donationMessage"
              >
                Message:
              </label>
              <textarea
                class="shadow appearance-none border rounded w-80 py-2 px-3 h-24 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
                id="donationMessage"
                type="text"
                :maxlength="msgMaxChar"
                v-model="donationMessage"
              ></textarea>
            </div>
            <div class="mb-6">
              <label class="block text-gray-500 font-bold">
                <input
                  class="mr-2 leading-tight accent-green-700"
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
                  aria-current="page"
                  :class="isCashAddrs ? unselectedTabClass : selectedTabClass"
                  @click="onQRChange"
                  >Address</a
                >
              </li>
              <div v-if="tknEnabled">
                <li class="me-2">
                  <a
                    :class="isCashAddrs ? selectedTabClass : unselectedTabClass"
                    @click="onQRChange"
                    >Cash Address</a
                  >
                </li>
              </div>
            </ul>

            <div class="w-80" v-show="isValidUser">
              <span class="text-xs">{{
                !isCashAddrs ? walletAddress : cashAddress
              }}</span>
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
              <p v-if="isValidSend" class="text-emerald-700 text-xs italic">
                Message Sent!
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
