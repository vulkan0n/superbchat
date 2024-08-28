import { createRouter, createWebHistory } from "vue-router";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      name: "home",
      component: () => import("../views/HomeView.vue"),
    },
    {
      path: "/user/login",
      name: "login",
      component: () => import("../views/LoginView.vue"),
    },
    {
      path: "/user/signup",
      name: "signup",
      component: () => import("../views/SignUpView.vue"),
    },
    {
      path: "/dashboard",
      name: "dashboard",
      component: () => import("../views/DashboardView.vue"),
    },
    {
      path: "/privacy-policy",
      name: "privacy-policy",
      component: () => import("../views/PrivacyPolicyView.vue"),
    },
    {
      path: "/:user",
      name: "superbchat",
      component: () => import("../views/SuperbchatView.vue"),
    },
    {
      path: "/404",
      name: "404",
      component: () => import("../views/404View.vue"),
    },
  ],
});

export default router;
