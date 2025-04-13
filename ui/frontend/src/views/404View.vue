<template>
  <main>
    <img
      src="../assets/404.png"
      alt="Error 404 - Not Found"
      class="rounded-2xl w-full max-w-[40rem] animate-fade-in"
    />
  </main>
</template>

<style scoped>
@keyframes fade-in {
  0% {
    opacity: 0;
    transform: scale(0.95);
  }
  100% {
    opacity: 1;
    transform: scale(1);
  }
}

.animate-fade-in {
  animation: fade-in 0.6s ease-out forwards;
}
</style>
