<script>
import { onMounted, onUnmounted, ref } from "vue";
import { BaseWallet, Wallet, BCMR, TokenSendRequest } from "mainnet-js";
import Token from "./Token.vue";

export default {
  components: { Token },
  setup() {
    const isCashAddrs = ref(false);
    const walletAddress = ref();
    const cashAddress = ref();
    const walletBalance = ref(0);
    const withdrawAddr = ref("");
    const withdrawTknAddr = ref("");
    const seedPhrase = ref("");

    const invalidAddrError = ref(false);
    const txSent = ref(false);
    const showSeedPhrase = ref(false);
    const tokenList = ref([]);
    const nftList = ref([]);

    let intervalId;
    let wallet;

    const errorClass = ref("border-red-500");
    const selectedTabClass = ref(
      "inline-block p-4 text-blue-600 bg-gray-100 rounded-t-lg active dark:bg-emerald-600 dark:text-white"
    );
    const unselectedTabClass = ref(
      "inline-block p-4 rounded-t-lg hover:text-emerald-300 hover:bg-gray-50 dark:text-gray-700 dark:hover:bg-emerald-800 dark:hover:text-gray-300"
    );

    function onQRChange() {
      getTokensInfo();
      const currentVal = isCashAddrs.value;
      isCashAddrs.value = !currentVal;
    }

    function onShowSeedPhrase() {
      const currentVal = showSeedPhrase.value;
      showSeedPhrase.value = !currentVal;
    }

    async function sendBCH() {
      try {
        if (withdrawAddr.value == "") throw "Empty Address";
        const txData = await wallet.sendMax(addBchPrefix(withdrawAddr.value));
        console.log(txData);
        await fetchBalance();
        txSent.value = true;
        setTimeout(() => {
          txSent.value = false;
        }, 3000);
      } catch (error) {
        console.error("Failed to send:", error);
        invalidAddrError.value = true;
        setTimeout(() => {
          invalidAddrError.value = false;
        }, 3000);
      }
    }

    async function sendTokens() {
      const tokensBalance = await wallet.getAllTokenBalances();
      const nftBalance = await wallet.getAllNftTokenBalances();

      let txs = [];
      for (const [key, value] of Object.entries(tokensBalance)) {
        const tknReq = new TokenSendRequest({
          cashaddr: addBchPrefix(withdrawTknAddr.value),
          amount: value,
          tokenId: key,
        });
        txs.push(tknReq);
      }
      for (const [key, value] of Object.entries(nftBalance)) {
        const tknReq = new TokenSendRequest({
          cashaddr: addBchPrefix(withdrawTknAddr.value),
          amount: value,
          tokenId: key,
        });
        txs.push(tknReq);
      }
      console.log(txs);
      try {
        if (withdrawTknAddr.value == "") throw "Empty Address";
        const txData = await wallet.send(txs);
        console.log(txData);
        getTokensInfo();
        await fetchBalance();
        txSent.value = true;
        setTimeout(() => {
          txSent.value = false;
        }, 3000);
      } catch (error) {
        console.error("Failed to send:", error);
        invalidAddrError.value = true;
        setTimeout(() => {
          invalidAddrError.value = false;
        }, 3000);
      }
    }

    function addBchPrefix(address) {
      if (!address.startsWith("bitcoincash:")) {
        return `bitcoincash:${address}`;
      }
      return address;
    }

    const fetchBalance = async () => {
      walletAddress.value = wallet.address;
      cashAddress.value = wallet.tokenaddr;
      walletBalance.value = await wallet.getBalance("bch");
    };

    const getTokensInfo = async () => {
      tokenList.value = [];
      nftList.value = [];
      const tokensBalance = await wallet.getAllTokenBalances();
      const nftBalance = await wallet.getAllNftTokenBalances();
      await pushResponse(tokensBalance, false);
      await pushResponse(nftBalance, true);
    };

    async function pushResponse(tokensBalance, isNFT) {
      for (const [key, value] of Object.entries(tokensBalance)) {
        await BCMR.addMetadataRegistryAuthChain({
          transactionHash: key,
          followToHead: false,
        });
        const tokenInfo = await wallet.getTokenInfo(key);
        const tknInfoObj = {
          tknCategoryId: key,
          tknSymbol: tokenInfo.token.symbol,
          tknName: tokenInfo.name,
          tknAmount: value,
          tknLogo: convertIpfsUrl(tokenInfo.uris.icon),
        };
        if (isNFT) nftList.value.push(tknInfoObj);
        else tokenList.value.push(tknInfoObj);
      }
    }

    function convertIpfsUrl(url) {
      if (url.startsWith("ipfs://")) {
        return url.replace("ipfs://", "https://ipfs.io/ipfs/");
      }
      return url;
    }

    onMounted(async () => {
      let user = "vulkan0n";
      BaseWallet.StorageProvider = IndexedDBProvider;
      wallet = await Wallet.named(`user:${user}`);
      seedPhrase.value = wallet.mnemonic;
      fetchBalance();
      intervalId = setInterval(fetchBalance, 5000);
    });

    onUnmounted(() => {
      clearInterval(intervalId);
    });

    return {
      errorClass,
      selectedTabClass,
      unselectedTabClass,
      walletAddress,
      walletBalance,
      isCashAddrs,
      cashAddress,
      onQRChange,
      withdrawAddr,
      withdrawTknAddr,
      invalidAddrError,
      txSent,
      sendBCH,
      sendTokens,
      tokenList,
      nftList,
      seedPhrase,
      showSeedPhrase,
      onShowSeedPhrase,
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
            <div>
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
                <li class="me-2">
                  <a
                    :class="isCashAddrs ? selectedTabClass : unselectedTabClass"
                    @click="onQRChange"
                    >Cash Address</a
                  >
                </li>
              </ul>

              <div class="w-80">
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
              <label
                class="text-gray-700 text-m font-bold mb-2 mr-2"
                for="alert-url"
                v-show="!isCashAddrs"
                >Balance: {{ walletBalance }} BCH</label
              >
              <label
                class="text-gray-700 text-m font-bold mb-2 mr-2"
                for="alert-url"
                v-show="isCashAddrs && tokenList.length > 0"
                >Token Balance:
              </label>
              <div
                v-for="tokenItem of tokenList"
                v-show="isCashAddrs"
                :key="tokenItem.tknCategoryId"
              >
                <token v-bind="tokenItem" />
              </div>
              <label
                class="text-gray-700 text-m font-bold mb-2 mr-2"
                for="alert-url"
                v-show="isCashAddrs && nftList.length > 0"
                >NFTs:
              </label>
              <div
                v-for="nftItem of nftList"
                v-show="isCashAddrs"
                :key="nftItem.tknCategoryId"
              >
                <token v-bind="nftItem" />
              </div>
            </div>
            <hr class="my-5" />
            <div class="space-y-2" v-show="!isCashAddrs">
              <label class="text-gray-700 text-m font-bold mb-2 mr-2"
                >Withdraw BCH Address</label
              >
              <small v-show="invalidAddrError" class="text-red-700 font-bold"
                ><label>Invalid address</label></small
              >
              <small v-show="txSent" class="text-green-700 font-bold"
                ><label>Sent!</label></small
              >
              <div class="flex space-x-2">
                <input
                  class="shadow border block rounded w-full mr-3 py-2 px-3 text-gray-700 leading-normal focus:outline-none focus:shadow-outline"
                  :class="
                    invalidAddrError
                      ? errorClass
                      : txSent
                      ? 'border-green-500'
                      : ''
                  "
                  type="text"
                  v-model="withdrawAddr"
                  placeholder="bitcoincash:qz..."
                />
                <button
                  class="bg-emerald-500 hover:bg-emerald-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
                  @click="sendBCH"
                >
                  Send
                </button>
              </div>
            </div>
            <div class="space-y-2" v-show="isCashAddrs">
              <label class="text-gray-700 text-m font-bold mb-2 mr-2"
                >Withdraw CashTokens Address</label
              >
              <small v-show="invalidAddrError" class="text-red-700 font-bold"
                ><label>Invalid address</label></small
              >
              <small v-show="txSent" class="text-green-700 font-bold"
                ><label>Sent!</label></small
              >
              <div class="flex space-x-2">
                <input
                  class="shadow border block rounded w-full mr-3 py-2 px-3 text-gray-700 leading-normal focus:outline-none focus:shadow-outline"
                  :class="
                    invalidAddrError
                      ? errorClass
                      : txSent
                      ? 'border-green-500'
                      : ''
                  "
                  type="text"
                  v-model="withdrawTknAddr"
                  placeholder="bitcoincash:zz..."
                />
                <button
                  class="bg-emerald-500 hover:bg-emerald-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
                  @click="sendTokens"
                >
                  Send
                </button>
              </div>
            </div>
            <hr class="my-5" />
            <div>
              <label class="text-gray-700 text-m font-bold mr-2"
                >Seed Phrase</label
              >
              <textarea
                readonly
                class="shadow border block rounded w-full h-24 mt-2 mr-3 py-2 px-3 text-gray-700 leading-normal focus:outline-none focus:shadow-outline"
                :class="showSeedPhrase ? 'blur-none' : 'blur-sm'"
                v-model="seedPhrase"
                @click="onShowSeedPhrase"
              ></textarea>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
