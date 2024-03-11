<script>
import { useRouter, useRoute } from "vue-router";
import { onMounted, ref } from "vue";
import { BaseWallet, Wallet } from "mainnet-js";

const fakeUsers = ["vulkan0n", "pepe"];
const walletAddress = ref();
const isValidUser = ref(true);

export default {
  setup() {
    const user = useRoute().params.user;
    const router = useRouter();

    onMounted(async () => {
      if (!fakeUsers.includes(user)) {
        router.push({ name: "404" });
        isValidUser.value = false;
      } else {
        BaseWallet.StorageProvider = IndexedDBProvider;
        const wallet = await Wallet.named(`user:${user}`);
        console.log(wallet);
        walletAddress.value = wallet.address;
      }
    });

    return { user, walletAddress, isValidUser };
  },
};
</script>

<template>
  <main v-show="isValidUser">
    <h1>Superbchat {{ user }}</h1>

    <h3>{{ walletAddress }}</h3>
    <qr-code
      :contents="walletAddress"
      module-color="#1c7d43"
      position-ring-color="#13532d"
      position-center-color="#70c559"
      style="background-color: #fff"
      class="qr"
    >
      <img src="../assets/bch.svg" slot="icon" />
    </qr-code>
  </main>
</template>
