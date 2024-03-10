<script>
import { useRouter, useRoute } from "vue-router";
import { onBeforeMount } from "vue";
import { onMounted, ref } from "vue";
import { BaseWallet, Wallet } from "mainnet-js";

const fakeUsers = ["vulkan0n", "pepe"];
const walletAddress = ref();

export default {
  setup() {
    const user = useRoute().params.user;
    const router = useRouter();

    onBeforeMount(() => {
      if (!fakeUsers.includes(user)) router.push({ name: "404" });
    });

    onMounted(async () => {
      BaseWallet.StorageProvider = IndexedDBProvider;
      const wallet = await Wallet.named(`user:${user}`);
      console.log(wallet);
      walletAddress.value = wallet.address;
    });

    return { user, walletAddress };
  },
};
</script>

<template>
  <main>
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
