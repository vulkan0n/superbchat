<script>
import { ref } from "vue";
import axios from "axios";
import bchLogo from "../assets/bitcoin-cash-bch-logo.svg";

export default {
  props: {
    id: {
      type: Number,
      required: true,
    },
    name: {
      type: String,
      required: true,
    },
    amount: {
      type: Number,
      required: true,
    },
    message: {
      type: String,
      required: true,
    },
    created: {
      type: String,
      required: true,
    },
    isCashToken: {
      type: Boolean,
      required: false,
    },
    tknCategory: {
      type: String,
      required: false,
    },
    tknSymbol: {
      type: String,
      required: false,
    },
    tknLogo: {
      type: String,
      required: false,
    },
    txId: {
      type: String,
      required: false,
    },
  },
  setup(props) {
    const isValidDelete = ref(false);

    async function onDelete() {
      try {
        const response = await axios.delete(`/superbchat/${props.id}`);

        if (response.status === 200) {
          isValidDelete.value = true; // Show success message

          setTimeout(() => {
            isValidDelete.value = false;
          }, 3000);
        }
      } catch (err) {
        console.error("Error deleting superchat:", err);
      } finally {
        window.location.reload();
      }
    }

    props.isCashToken = props.tknSymbol != "";
    if (!props.isCashToken) {
      props.tknSymbol = "BCH";
    }

    const createdDate = new Date(props.created);
    const createdDateStr = createdDate.toLocaleDateString();

    return { onDelete, createdDateStr, isValidDelete, bchLogo };
  },
};
</script>

<template>
  <article
    :key="name"
    class="py-3 px-6 my-2 max-w-md ring-black text-base ring-2 bg-white border-t border-gray-900 dark:border-gray-900 dark:bg-gray-500 rounded-xl"
  >
    <footer class="flex justify-between items-center mb-2">
      <div class="flex items-center">
        <a
          class="inline-flex items-center mr-3 text-sm text-gray-900 hover:underline dark:text-white font-semibold"
          :href="`https://blockchair.com/bitcoin-cash/transaction/${txId}`"
          target="_blank"
        >
          <img
            class="mr-2 w-6 h-6 rounded-full"
            :src="isCashToken ? tknLogo : bchLogo"
            alt="Token Logo"
          />{{ name }} ({{ amount }} {{ tknSymbol }})
        </a>
        <p class="text-sm mr-4 text-gray-600 dark:text-gray-400">
          <time pubdate="" :title="createdDateStr">{{ createdDateStr }}</time>
        </p>
      </div>
      <button
        id="dropdownComment4Button"
        class="inline-flex items-center p-2 text-sm font-medium text-center text-gray-500 dark:text-gray-40 bg-white rounded-lg hover:bg-gray-100 focus:ring-4 focus:outline-none focus:ring-gray-50 dark:bg-gray-700 dark:hover:bg-gray-800 dark:focus:ring-gray-600"
        type="button"
        @click="onDelete"
      >
        <svg
          width="16px"
          height="16px"
          viewBox="0 0 512 512"
          version="1.1"
          xmlns="http://www.w3.org/2000/svg"
          xmlns:xlink="http://www.w3.org/1999/xlink"
        >
          <title>Delete</title>
          <g
            id="Page-1"
            stroke="none"
            stroke-width="1"
            fill="none"
            fill-rule="evenodd"
          >
            <g
              id="trashcan"
              fill="currentColor"
              transform="translate(64.000000, 42.666667)"
            >
              <path
                d="M256,42.6666667 L128,42.6666667 L128,7.10542736e-15 L256,7.10542736e-15 L256,42.6666667 Z M170.666667,170.666667 L128,170.666667 L128,341.333333 L170.666667,341.333333 L170.666667,170.666667 Z M256,170.666667 L213.333333,170.666667 L213.333333,341.333333 L256,341.333333 L256,170.666667 Z M384,85.3333333 L384,128 L341.333333,128 L341.333333,426.666667 L42.6666667,426.666667 L42.6666667,128 L0,128 L0,85.3333333 L384,85.3333333 Z"
                id="Shape"
              ></path>
            </g>
          </g>
        </svg>
      </button>
    </footer>
    <p class="text-gray-500 dark:text-gray-900">
      {{ message }}
    </p>
  </article>
</template>
