<script>
import { onMounted, ref } from "vue";
import { RouterLink, useRouter } from "vue-router";
import axios from "axios";
import { Wallet } from "mainnet-js";

export default {
  setup() {
    const address = ref("bitcoincash:addressTest");
    const tknEnabled = ref(true);
    const tknAddress = ref("bitcoincash:tknAddressTest");
    const minDonation = ref(0.00000547);
    const defaultShowAmount = ref(true);
    const msgMaxChar = ref(300);

    const minDonationLowError = ref(false);
    const invalidAddressError = ref(false);
    const invalidTknAddressError = ref(false);

    const errorClass = ref("border-red-500");
    const isValidUpdate = ref(false);

    const router = useRouter();

    async function isValidAddress(address) {
      try {
        await Wallet.watchOnly(address);
        return true;
      } catch (error) {
        return false;
      }
    }

    async function verifyAndUpdateSettings() {
      minDonationLowError.value = minDonation.value < 0.00000547;

      invalidAddressError.value = !(await isValidAddress(address.value));

      invalidTknAddressError.value =
        tknEnabled.value &&
        (tknAddress.value === "" ||
          !(
            tknAddress.value.startsWith("bitcoincash:z") ||
            tknAddress.value.startsWith("z")
          ) ||
          !(await isValidAddress(tknAddress.value)));

      if (
        !minDonationLowError.value &&
        !invalidAddressError.value &&
        !invalidTknAddressError.value
      ) {
        await postSettings();
      }
    }

    async function postSettings() {
      try {
        const data = {
          address: address.value,
          tknsEnabled: tknEnabled.value,
          tknAddress: tknAddress.value,
          messageMaxChars: Number(msgMaxChar.value),
          minDonation: minDonation.value,
          showAmount: defaultShowAmount.value,
        };

        console.log(data);

        const response = await axios.post("/settings", data, {
          headers: {
            Authorization: `Bearer ${localStorage.getItem("token")}`,
            "Content-Type": "application/json",
          },
        });

        if (response.statusText == "OK") {
          console.log(response);
          isValidUpdate.value = true;
          setTimeout(() => {
            isValidUpdate.value = false;
          }, 3000);
        } else {
          console.log(response);
          router.push("/login");
        }
      } catch (err) {
        console.log(err);
        router.push("/login");
      }
    }

    onMounted(async () => {
      try {
        const settingsInfoResponse = await axios.get(
          "/user-id/" + localStorage.getItem("userId")
        );
        if (settingsInfoResponse.statusText == "OK") {
          address.value = settingsInfoResponse.data.address;
          tknEnabled.value = settingsInfoResponse.data.tknsEnabled;
          tknAddress.value = settingsInfoResponse.data.tknAddress;
          minDonation.value = settingsInfoResponse.data.minDonation;
          defaultShowAmount.value = settingsInfoResponse.data.showAmount;
          msgMaxChar.value = settingsInfoResponse.data.messageMaxChars;
        } else {
          router.push({ name: "404" });
        }
      } catch (err) {
        router.push({ name: "404" });
      }
    });

    return {
      RouterLink,
      address,
      tknEnabled,
      tknAddress,
      minDonation,
      defaultShowAmount,
      msgMaxChar,
      minDonationLowError,
      errorClass,
      verifyAndUpdateSettings,
      isValidUpdate,
      invalidAddressError,
      invalidTknAddressError,
    };
  },
};
</script>

<template>
  <div class="flex flex-col justify-center items-center">
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
                :class="invalidAddressError ? errorClass : ''"
                id="address"
                type="text"
                placeholder="bitcoincash:"
                v-model="address"
              />
              <p v-if="invalidAddressError" class="text-red-500 text-xs italic">
                Please choose a valid BCH address.
              </p>
            </div>
            <div class="mb-4" v-if="tknEnabled">
              <label
                class="block text-gray-700 text-sm font-bold mb-2"
                for="address"
              >
                CashTokens Address
              </label>
              <input
                class="shadow appearance-none border rounded w-80 py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
                :class="invalidTknAddressError ? errorClass : ''"
                id="tknAddress"
                type="text"
                placeholder="bitcoincash:"
                v-model="tknAddress"
              />
              <p v-if="invalidTknAddressError" class="text-red-500 text-xs italic">
                Please choose a valid CashTokens address.
              </p>
            </div>
            <div class="mb-4">
              <label
                class="block text-gray-700 text-sm font-bold mb-2"
                for="minDonation"
              >
                Minimum Donation Amount
              </label>
              <input
                class="shadow appearance-none border rounded w-80 py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
                :class="minDonationLowError ? errorClass : ''"
                id="minDonation"
                type="number"
                step="0.000001"
                v-model="minDonation"
              />
              <p v-if="minDonationLowError" class="text-red-500 text-xs italic">
                Minimun donation amount must be higher than 0.00000546
              </p>
            </div>
            <div class="mb-4">
              <label
                class="block text-gray-700 text-sm font-bold mb-2"
                for="minDonation"
              >
                Max Message Characters Lenght
              </label>
              <select
                id="size-select"
                v-model="msgMaxChar"
                class="shadow border rounded leading-tight focus:outline-none focus:shadow-outline block w-full p-2.5"
              >
                <option value="150" :selected="msgMaxChar === 150">
                  Small (150 Characters)
                </option>
                <option value="300" :selected="msgMaxChar === 300">
                  Medium (300 Characters)
                </option>
                <option value="500" :selected="msgMaxChar === 500">
                  Large (500 Characters)
                </option>
              </select>
            </div>
            <div class="mb-6">
              <label class="block text-gray-500 font-bold">
                <input
                  class="mr-2 leading-tight accent-emerald-700"
                  type="checkbox"
                  v-model="defaultShowAmount"
                />
                <span class="text-sm"> Default "Show Amount" checked </span>
              </label>
            </div>
            <div class="mb-6">
              <label class="block text-gray-500 font-bold">
                <input
                  class="mr-2 leading-tight accent-emerald-700"
                  type="checkbox"
                  v-model="tknEnabled"
                />
                <span class="text-sm"> CashTokens enabled </span>
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
              <p v-if="isValidUpdate" class="text-emerald-700 text-xs italic">
                Update successful.
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
