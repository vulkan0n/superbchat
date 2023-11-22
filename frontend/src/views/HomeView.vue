<script>
import { onMounted, ref } from "vue";
import { Wallet } from "mainnet-js";

const walletAddress = ref();

export default {
  setup() {
    onMounted(async () => {
      const wallet = await Wallet.newRandom();
      walletAddress.value = wallet.address;
      console.log(walletAddress.value);
    });
    return { walletAddress };
  },
};
</script>

<template>
  <div class="wrapper">
    <nav>
      <RouterLink to="/user/login">Login</RouterLink>
      <RouterLink to="/user/signup">Signup</RouterLink>
    </nav>
  </div>
  <main>
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

<style scoped>
.qr {
  width: 25rem;
}
</style>
