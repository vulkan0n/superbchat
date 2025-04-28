<template>
  <div
    class="bg-transparent fixed inset-0 flex items-center justify-center pointer-events-none"
  >
    <!-- Donation Alert -->
    <div
      v-if="alert"
      class="p-5 bg-white rounded-lg shadow-xl border-l-8 border-green-500 animate-pop-in max-w-md w-full"
    >
      <div class="flex items-center space-x-4">
        <img
          :src="alert.isTkn ? alert.tknLogo : bchLogo"
          class="h-24 w-24 rounded-full object-cover"
          alt="Avatar"
        />
        <div class="space-y-1">
          <div class="text-2xl font-bold text-gray-900 leading-tight">
            {{ alert.username }}
          </div>
          <div class="text-base text-gray-700 leading-tight">donated</div>
          <div
            v-if="alert.showAmount"
            class="text-3xl font-extrabold text-green-600 leading-tight"
          >
            {{ alert.amount }} {{ alert.tknSymbol }}
          </div>
          <div class="text-lg text-gray-800 mt-2 leading-snug">
            {{ alert.message }}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from "vue";
import { useRoute } from "vue-router";
import bchLogo from "../assets/bch.svg";

const alert = ref(null);

function triggerDonationAlert(
  username,
  amount,
  message,
  isTkn,
  showAmount,
  tknSymbol,
  tknLogo
) {
  alert.value = {
    username,
    amount,
    message,
    isTkn,
    showAmount,
    tknSymbol: isTkn ? tknSymbol : "BCH",
    tknLogo,
  };

  // Clear the alert after 8 seconds
  setTimeout(() => {
    alert.value = null;
  }, 8000);
}

const widgetId = useRoute().params.uuid;

let socket = null;
const wsUrl = window.location.origin.replace(/^http/, "ws") + "/ws";

const connectWebSocket = () => {
  socket = new WebSocket(wsUrl);

  socket.onmessage = (event) => {
    const chat = JSON.parse(event.data);
    if (chat.widget_id === widgetId) {
      triggerDonationAlert(
        chat.name,
        chat.amount,
        chat.message,
        chat.isTkn,
        !chat.isHidden,
        chat.tknSymbol,
        chat.tknLogo
      );
    }
  };

  socket.onerror = (error) => {
    console.error("WebSocket Error:", error);
  };
};

onMounted(() => {
  connectWebSocket();
});

onUnmounted(() => {
  if (socket) {
    socket.close();
  }
});
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
