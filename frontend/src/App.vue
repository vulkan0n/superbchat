<script>
import { RouterLink, RouterView } from "vue-router";
export default {
  setup() {
    return { RouterLink, RouterView };
  },
};
</script>

<template>
  <header>
    <RouterLink to="/">
      <div class="title">
        <img
          alt="BCH logo"
          class="logo"
          src="@/assets/bitcoin-cash-bch-logo.svg"
          width="45"
          height="45"
        />
        <h1>Superbchat!</h1>
      </div>
    </RouterLink>
  </header>
  <RouterView />
</template>

<style scoped>
.title {
  display: flex;
  place-items: center;
}

header {
  line-height: 1.5;
  max-height: 100vh;
}

.logo {
  display: block;
  margin: 0 1rem 0 0;
}

nav {
  width: 100%;
  font-size: 12px;
  text-align: center;
  margin-top: 2rem;
}

nav a.router-link-exact-active {
  color: var(--color-text);
}

nav a.router-link-exact-active:hover {
  background-color: transparent;
}

nav a {
  display: inline-block;
  padding: 0 1rem;
  border-left: 1px solid var(--color-border);
}

nav a:first-of-type {
  border: 0;
}

@media (min-width: 1024px) {
  header {
    padding-right: calc(var(--section-gap) / 2);
  }

  .logo {
    margin: 0 1rem 0 0;
  }

  header .wrapper {
    display: flex;
    place-items: flex-start;
    flex-wrap: wrap;
  }

  nav {
    text-align: left;
    margin-left: -1rem;
    font-size: 1rem;

    padding: 1rem 0;
    margin-top: 1rem;
  }
}
</style>
