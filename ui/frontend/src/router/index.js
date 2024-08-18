import { createRouter, createWebHistory } from "vue-router";
import HomeView from "../views/HomeView.vue";
import LoginView from "../views/LoginView.vue";
import DashboardView from "../views/DashboardView.vue";
import SuperbchatView from "../views/SuperbchatView.vue";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      name: "home",
      component: HomeView,
    },
    {
      path: "/user/login",
      name: "login",
      component: LoginView,
    },
    {
      path: "/user/signup",
      name: "signup",
      component: () => import("../views/SignUpView.vue"),
    },
    {
      path: "/dashboard",
      name: "dashboard",
      component: DashboardView,
    },
    {
      path: "/privacy-policy",
      name: "privacy-policy",
      component: () => import("../views/PrivacyPolicyView.vue"),
    },
    {
      path: "/:user",
      name: "superbchat",
      component: SuperbchatView,
    },
    {
      path: "/404",
      name: "404",
      component: () => import("../views/404View.vue"),
    },
  ],
});

export default router;
