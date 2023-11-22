<script>
import { useRouter, useRoute } from "vue-router";
import { onBeforeMount } from "vue";

const fakeUsers = ["vulkan0n", "pepe"];
export default {
  setup() {
    const user = useRoute().params.user;
    const router = useRouter();

    onBeforeMount(() => {
      if (!fakeUsers.includes(user)) router.push({ name: "404" });
    });

    return { user };
  },
};
</script>

<template>
  <main>
    <h1>Superbchat {{ user }}</h1>
  </main>
</template>
