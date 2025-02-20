<script>
import QRCodeVue3 from "qrcode-vue3";
import { useRouter, useRoute } from "vue-router";
import { onMounted, ref } from "vue";
import { Wallet, BCMR } from "mainnet-js";
import axios from "axios";

const walletAddress = ref();
const tknEnabled = ref(true);
const tknAddress = ref();
const isValidUser = ref(true);
const minDonation = ref(0.00000547);
const donationAmount = ref(0.00000547);
const msgMaxChar = ref(300);
const showAmount = ref(true);
const waitingForTx = ref(false);
let wallet;
let tknWallet;

export default {
  components: { QRCodeVue3 },
  setup() {
    const donatorName = ref("Anonymous");
    const donationMessage = ref("");
    const isTknAddrs = ref(false);
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
      if (!isTknAddrs.value) {
        await awaitDontationBCH();
      } else {
        await awaitDontationTKN();
      }
    }

    let txId;
    let tknCategoryId;
    let tknSymbol;
    let tknLogo;


    async function awaitDontationBCH() {
      if (donationAmount.value < minDonation.value) {
        donationLowError.value = true;
        setTimeout(() => {
          donationLowError.value = false;
        }, 3000);
      }
      if (!donationLowError.value) {
        waitingForTx.value = true;
        const cancelWatch = wallet.watchAddressTransactions(async (tx) => {
          if (tx.vout[0].value == donationAmount.value) {
            txId = tx.txid;
            tknCategoryId = "";
            waitingForTx.value = false;
            tknSymbol = "";
            tknLogo = "";
            await sendSuperbchat();
            await cancelWatch();
          }
        });
      }
    }

    async function awaitDontationTKN() {
      waitingForTx.value = true;
      tknSymbol = "";
      tknLogo = "";

      const cancelWatch = tknWallet.watchAddressTransactions(async (tx) => {
        txId = tx.txid;
        tknCategoryId = tx.vout[0].tokenData.category;
        try {
          await BCMR.addMetadataRegistryAuthChain({
            transactionHash: tknCategoryId,
          });
          const info = BCMR.getTokenInfo(tknCategoryId);
          tknSymbol = info.token.symbol;
          tknLogo = info.uris.icon;
          const rawAmount = parseFloat(tx.vout[0].tokenData.amount);
          donationAmount.value = rawAmount / Math.pow(10, info.token.decimals);
        } catch (err) {
          console.warn("BCMR error with categoryId: " + tknCategoryId);
          console.warn(err);
          tknSymbol = "CASHTOKEN";
          tknLogo = "https://cashonize.com/images/tokenicon.png";
          donationAmount.value = parseFloat(tx.vout[0].tokenData.amount);
        }
        waitingForTx.value = false;
        await sendSuperbchat();
        await cancelWatch();
      });
    }

    async function sendSuperbchat() {
      try {
        let postBody = {
          txId: txId,
          name: donatorName.value,
          message: donationMessage.value,
          amount: donationAmount.value,
          tknCategory: tknCategoryId,
          tknSymbol: tknSymbol,
          tknLogo: tknLogo,
          isHidden: !showAmount.value,
          recipient: userId,
          isTkn: isTknAddrs.value,
        };
        const superbchatResponse = await axios.post("/superbchat", postBody);
        if (superbchatResponse.status === 200) {
          console.log("Superchat Sent");
          isValidSend.value = true;
          donationMessage.value = "";
          setTimeout(() => {
            isValidSend.value = false;
          }, 3000);
        }
      } catch (err) {
        console.error(err);
      } finally {
        setTimeout(() => {
          router.push(route.fullPath);
        }, 3000);
      }
    }

    function onQRChange() {
      const currentVal = isTknAddrs.value;
      isTknAddrs.value = !currentVal;
    }

    onMounted(async () => {
      try {
        const userInfoResponse = await axios.get("/user/" + user);
        if (userInfoResponse.statusText == "OK") {
          walletAddress.value = userInfoResponse.data.address;
          tknEnabled.value = userInfoResponse.data.tknsEnabled;
          tknAddress.value = userInfoResponse.data.tknAddress;
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
      if (walletAddress.value) {
        wallet = await Wallet.watchOnly(walletAddress.value);
      }
      if (tknEnabled.value && tknAddress.value) {
        tknWallet = await Wallet.watchOnly(tknAddress.value);
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
      isTknAddrs,
      tknAddress,
      onQRChange,
      tknEnabled,
      minDonation,
      msgMaxChar,
      isValidSend,
      waitingForTx,
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
            <div v-if="!isTknAddrs" class="mb-4">
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
                  :class="isTknAddrs ? unselectedTabClass : selectedTabClass"
                  @click="onQRChange"
                  >Address</a
                >
              </li>
              <div v-if="tknEnabled">
                <li class="me-2">
                  <a
                    :class="isTknAddrs ? selectedTabClass : unselectedTabClass"
                    @click="onQRChange"
                    >Token Address</a
                  >
                </li>
              </div>
            </ul>

            <div
              class="w-80"
              v-show="isValidUser"
              :class="waitingForTx ? 'blur-none' : 'blur-sm'"
            >
              <span class="text-xs">{{
                !isTknAddrs ? walletAddress : tknAddress
              }}</span>
              <div class="m-4">
                <QRCodeVue3
                  :value="
                    !isTknAddrs
                      ? walletAddress + '?amount=' + donationAmount
                      : tknAddress
                  "
                  :key="
                    !isTknAddrs
                      ? walletAddress + '?amount=' + donationAmount
                      : tknAddress
                  "
                  :dots-options="{ type: 'dots', color: '#1c7d43' }"
                  :corners-square-options="{
                    type: 'extra-rounded',
                    color: '#13532d',
                  }"
                  :corners-dot-options="{
                    type: 'dot',
                    color: '#70c559',
                  }"
                />
              </div>
            </div>
            <div class="flex items-center justify-between">
              <button
                class="bg-green-600 hover:bg-green-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
                type="button"
                @click="verifyAndPay"
                :disabled="waitingForTx"
              >
                {{
                  !waitingForTx
                    ? "Prepare donation"
                    : "Waiting for transaction..."
                }}
              </button>
              <div v-if="waitingForTx" class="ml-2">
                <svg
                  class="animate-spin h-5 w-5 text-green-600"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                >
                  <circle cx="12" cy="12" r="10"></circle>
                  <path d="M12 2v4"></path>
                  <path d="M12 18v4"></path>
                  <path d="M4.93 4.93l2.83 2.83"></path>
                  <path d="M16.24 16.24l2.83 2.83"></path>
                  <path d="M2 12h4"></path>
                  <path d="M18 12h4"></path>
                  <path d="M4.93 19.07l2.83-2.83"></path>
                  <path d="M16.24 7.76l2.83-2.83"></path>
                </svg>
              </div>
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
