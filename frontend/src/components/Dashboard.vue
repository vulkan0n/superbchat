<script>
const fakeAlertUrl =
  "192.168.100.6:8900/alert/81f1cd27-d510-4579-97c8-9613e5f63fb2";

const fakeMessages = [
  {
    sender: "Vulkan0n",
    amount: 0.0001,
    message: "Cool! Let´s go BCH!",
  },
  {
    sender: "Anonymous",
    amount: 0.00014,
    message: "Awesome!",
  },
  {
    sender: "VknToken",
    amount: 0,
    message: "Get some tokens!",
    isCashToken: true,
    cashToken: {
      symbol: "VKN",
      amount: 10,
      logo: "https://gist.githubusercontent.com/vulkan0n/74922802a5d3a8861765e882c3a2db1a/raw/logo.png",
    },
  },
];

export default {
  setup() {
    return { fakeMessages, fakeAlertUrl };
  },
};
</script>

<template>
  <div>
    <h1>{{ fakeAlertUrl }}</h1>
    <li
      v-for="superbchat in fakeMessages"
      :key="superbchat.sender"
      class="my-0.5"
    >
      {{ superbchat.sender }} : {{ superbchat.message }}
    </li>
  </div>
</template>
