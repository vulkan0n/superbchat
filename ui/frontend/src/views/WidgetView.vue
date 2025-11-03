<template>
  <div
    class="bg-white fixed inset-0 flex items-center justify-center pointer-events-none"
  >
    <!-- Donation Alert -->
    <div
      v-if="alert"
      class="p-6 bg-white rounded-lg shadow-xl border-l-8 border-green-500 animate-pop-in"
    >
      <div class="flex items-center space-x-4">
        <img
          v-if="!alert.isTkn"
          src="../assets/bch.svg"
          class="h-12 w-12 rounded-full"
          alt="Avatar"
        /><img
          v-if="alert.isTkn"
          src="https://gist.github.com/vulkan0n/74922802a5d3a8861765e882c3a2db1a/raw/1c8b5bbcbb8fa282c672fbcebefbf4ac52bfdf34/logo.png"
          class="h-12 w-12 rounded-full"
          alt="Avatar"
        />
        <div>
          <div class="text-lg font-bold text-gray-900">
            {{ alert.username }}
          </div>
          <div class="text-sm text-gray-600">donated</div>
          <div v-if="alert.showAmount" class="text-2xl font-bold text-green-600">
            BCH {{ alert.amount }}
          </div>
          <div class="text-sm text-gray-500">{{ alert.message }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from "vue";

const alert = ref(null);

function triggerDonationAlert(username, amount, message, isTkn, showAmount) {
  alert.value = {
    username,
    amount,
    message,
    isTkn,
    showAmount,
  };

  // Clear the alert after 5 seconds
  setTimeout(() => {
    alert.value = null;
  }, 5000);
}

triggerDonationAlert("Milton", 2.5, "Let's go BCH!", true, false);
setTimeout(() => {
triggerDonationAlert("Rojer", 0.005, "Excelent job", false, true);
  }, 5000);

</script>

<style>
/* Custom pop-in animation */
@keyframes popIn {
  from {
    transform: scale(0.5);
    opacity: 0;
  }
  to {
    transform: scale(1);
    opacity: 1;
  }
}

.animate-pop-in {
  animation: popIn 0.5s ease-out;
}
</style>
