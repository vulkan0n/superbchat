<script>
import Message from "./Message.vue";
import { ref, onMounted } from "vue";
import axios from "axios";
const fakeAlertUrlObj =
  "192.168.100.6:8900/alert/81f1cd27-d510-4579-97c8-9613e5f63fb2";

export default {
  components: { Message },
  setup() {
    const Messages = ref("");
    const fakeAlertUrl = ref("");
    const copiedUrl = ref(false);

    const unsecuredCopyToClipboard = (text) => {
      const textArea = document.createElement("textarea");
      textArea.value = text;
      document.body.appendChild(textArea);
      textArea.focus({ preventScroll: true });
      textArea.select();
      try {
        document.execCommand("copy");
      } catch (err) {
        console.error("Unable to copy to clipboard", err);
      }
      document.body.removeChild(textArea);
    };

    const copyToClipboard = () => {
      var content = fakeAlertUrl.value;
      if (window.isSecureContext && navigator.clipboard) {
        navigator.clipboard.writeText(content);
      } else {
        unsecuredCopyToClipboard(content);
      }
      copiedUrl.value = true;
      setTimeout(() => {
        copiedUrl.value = false;
      }, 2000);
    };
    onMounted(async () => {
      fakeAlertUrl.value = fakeAlertUrlObj;

      const token = localStorage.getItem("token");
      try {
        const superchatsResponse = await axios.post("/superbchat-get", {
          token,
        });
        if (superchatsResponse.statusText == "OK") {
          Messages.value = superchatsResponse.data;
          //fakeMessages.value = fakeMessagesObj;
        } else {
          console.log(superchatsResponse);
        }
      } catch (err) {
        console.log(err);
      }
    });

    return { Messages, fakeAlertUrl, copyToClipboard, copiedUrl };
  },
};
</script>

<template>
  <div class="flex md:flex-col justify-center">
    <div class="max-w-full mx-4 text-gray-600 items-center">
      <transition name="fade">
        <div
          v-if="copiedUrl"
          class="flex z-20 absolute left-1/3 md:left-1/2 md:-ml-20 -mt-4 rounded bg-emerald-500 text-white text-sm font-bold px-4 py-3"
          role="alert"
        >
          <svg
            class="fill-current w-4 h-4 mr-2"
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 20 20"
          >
            <path
              d="M12.432 0c1.34 0 2.01.912 2.01 1.957 0 1.305-1.164 2.512-2.679 2.512-1.269 0-2.009-.75-1.974-1.99C9.789 1.436 10.67 0 12.432 0zM8.309 20c-1.058 0-1.833-.652-1.093-3.524l1.214-5.092c.211-.814.246-1.141 0-1.141-.317 0-1.689.562-2.502 1.117l-.528-.88c2.572-2.186 5.531-3.467 6.801-3.467 1.057 0 1.233 1.273.705 3.23l-1.391 5.352c-.246.945-.141 1.271.106 1.271.317 0 1.357-.392 2.379-1.207l.6.814C12.098 19.02 9.365 20 8.309 20z"
            />
          </svg>
          <p>URL copied!</p>
        </div>
      </transition>
      <div
        id="mainBoard"
        class="rounded-lg shadow-lg p-3 md:p-8 w-full md:w-max m-auto bg-white overflow-auto overflow-ellipsis"
      >
        <div class="space-y-2">
          <label
            class="text-gray-700 text-m font-bold mb-2 mr-2"
            for="alert-url"
            >Widget URL</label
          >
          <small
            ><label
              >(Use this URL in OBS Studio, do not share it!)</label
            ></small
          >
          <div class="flex space-x-2">
            <input
              class="shadow border block rounded w-full mr-3 py-2 px-3 text-gray-700 leading-normal focus:outline-none focus:shadow-outline"
              id="fakeAlertUrl"
              type="text"
              :value="fakeAlertUrl"
              readonly
            />
            <button
              class="bg-emerald-500 hover:bg-emerald-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
              @click="copyToClipboard"
            >
              Copy
            </button>
          </div>
        </div>
        <hr class="my-5" />
        <div>
          <label class="text-gray-700 text-m font-bold mb-2">Messages</label>
          <div v-for="superbchat in Messages" :key="superbchat.seder">
            <message v-bind="superbchat" />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.5s;
}
.fade-enter, .fade-leave-to /* .fade-leave-active in <2.1.8 */ {
  opacity: 0;
}
</style>
