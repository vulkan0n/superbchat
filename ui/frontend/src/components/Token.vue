<script>
export default {
  props: {
    tknCategoryId: {
      type: String,
      required: true,
    },
    tknSymbol: {
      type: String,
      required: true,
    },
    tknName: {
      type: String,
      required: true,
    },
    tknAmount: {
      type: BigInt,
      required: true,
    },
    tknLogo: {
      type: String,
      required: true,
    },
  },
  setup(props) {
    return {};
  },
};
</script>

<template>
  <article
    :key="tknCategoryId"
    class="py-2 px-6 my-2 max-w-md ring-black text-base ring-2 bg-white border-t border-gray-900 dark:border-gray-900 dark:bg-gray-500 rounded-xl"
  >
    <div class="flex items-center">
      <img
        class="mr-7 w-10 h-10 border border-black rounded-full"
        :src="tknLogo"
        alt="Token Logo"
      />
      <div>
        <p
          class="inline-flex items-center mr-3 text-sm text-gray-900 dark:text-white font-semibold"
        >
          {{ tknAmount }} {{ tknSymbol }}
        </p>
        <p class="text-gray-500 dark:text-gray-900">
          {{ tknName }}
        </p>
      </div>
    </div>
  </article>
</template>
